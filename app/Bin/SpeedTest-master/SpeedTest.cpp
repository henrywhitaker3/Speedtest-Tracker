//
// Created by Francesco Laurita on 5/29/16.
//

#include <cmath>
#include <iomanip>
#include "SpeedTest.h"
#include "MD5Util.h"
#include <netdb.h>

SpeedTest::SpeedTest(float minServerVersion):
        mLatency(0),
        mUploadSpeed(0),
        mDownloadSpeed(0) {
    curl_global_init(CURL_GLOBAL_DEFAULT);
    mIpInfo = IPInfo();
    mServerList = std::vector<ServerInfo>();
    mMinSupportedServer = minServerVersion;
}

SpeedTest::~SpeedTest() {
    curl_global_cleanup();
    mServerList.clear();
}

bool SpeedTest::ipInfo(IPInfo& info) {

    if (!mIpInfo.ip_address.empty()){
        info = mIpInfo;
        return true;
    }

    std::stringstream oss;
    auto code = httpGet(SPEED_TEST_IP_INFO_API_URL, oss);
    if (code == CURLE_OK){
        auto values = SpeedTest::parseQueryString(oss.str());
        mIpInfo = IPInfo();
        mIpInfo.ip_address = values["ip_address"];
        mIpInfo.isp = values["isp"];
        mIpInfo.lat = std::stof(values["lat"]);
        mIpInfo.lon = std::stof(values["lon"]);
        values.clear();
        oss.clear();
        info = mIpInfo;
        return true;
    }

    return false;
}

const std::vector<ServerInfo>& SpeedTest::serverList() {
    if (!mServerList.empty())
        return mServerList;

    int http_code = 0;
    if (fetchServers(SPEED_TEST_SERVER_LIST_URL, mServerList, http_code) && http_code == 200){
        return mServerList;
    }
    return mServerList;
}


const ServerInfo SpeedTest::bestServer(const int sample_size, std::function<void(bool)> cb) {
    auto best = findBestServerWithin(serverList(), mLatency, sample_size, cb);
    SpeedTestClient client = SpeedTestClient(best);
    testLatency(client, SPEED_TEST_LATENCY_SAMPLE_SIZE, mLatency);
    client.close();
    return best;
}

bool SpeedTest::setServer(ServerInfo& server){
    SpeedTestClient client = SpeedTestClient(server);
    if (client.connect() && client.version() >= mMinSupportedServer){
        if (!testLatency(client, SPEED_TEST_LATENCY_SAMPLE_SIZE, mLatency)){
            return false;
        }
    } else {
        client.close();
        return false;
    }
    client.close();
    return true;

}

bool SpeedTest::downloadSpeed(const ServerInfo &server, const TestConfig &config, double& result, std::function<void(bool)> cb) {
    opFn pfunc = &SpeedTestClient::download;
    mDownloadSpeed = execute(server, config, pfunc, cb);
    result = mDownloadSpeed;
    return true;
}

bool SpeedTest::uploadSpeed(const ServerInfo &server, const TestConfig &config, double& result, std::function<void(bool)> cb) {
    opFn pfunc = &SpeedTestClient::upload;
    mUploadSpeed = execute(server, config, pfunc, cb);
    result = mUploadSpeed;
    return true;
}

const long &SpeedTest::latency() {
    return mLatency;
}

bool SpeedTest::jitter(const ServerInfo &server, long& result, const int sample) {
    auto client = SpeedTestClient(server);
    double current_jitter = 0;
    long previous_ms =  LONG_MAX;
    if (client.connect()){
        for (int i = 0; i < sample; i++){
            long ms = 0;
            if (client.ping(ms)){
                if (previous_ms == LONG_MAX) {
                    previous_ms = ms;
                } else {
                    current_jitter += std::abs(previous_ms - ms);
                }
            }
        }
        client.close();
    } else {
        return false;
    }

    result = (long) std::floor(current_jitter / sample);
    return true;
}


bool SpeedTest::share(const ServerInfo& server, std::string& image_url) {
    std::stringstream hash;
    hash << std::setprecision(0) << std::fixed << mLatency
    << "-" << std::setprecision(2) << std::fixed << (mUploadSpeed * 1000)
    << "-" << std::setprecision(2) << std::fixed << (mDownloadSpeed * 1000)
    << "-" << SPEED_TEST_API_KEY;
    std::string hex_digest = MD5Util::hexDigest(hash.str());
    std::stringstream post_data;
    post_data << "download=" << std::setprecision(2) << std::fixed << (mDownloadSpeed * 1000) << "&";
    post_data << "ping=" << std::setprecision(0) << std::fixed << mLatency << "&";
    post_data << "upload=" << std::setprecision(2) << std::fixed << (mUploadSpeed * 1000) << "&";
    post_data << "pingselect=1&";
    post_data << "recommendedserverid=" << server.id << "&";
    post_data << "accuracy=1&";
    post_data << "serverid=" << server.id << "&";
    post_data << "hash=";
    post_data << hex_digest;

    std::stringstream result;
    CURL *c = curl_easy_init();
    curl_easy_setopt(c, CURLOPT_REFERER, SPEED_TEST_API_REFERER);
    auto cres = SpeedTest::httpPost(SPEED_TEST_API_URL, post_data.str(), result, c);
    long http_code = 0;
    image_url.clear();
    if (cres == CURLE_OK){
        curl_easy_getinfo(c, CURLINFO_HTTP_CODE, &http_code);
        if (http_code == 200 && !result.str().empty()){
            auto data = SpeedTest::parseQueryString(result.str());
            if (data.count("resultid") == 1){
                image_url = "http://www.speedtest.net/result/" + data["resultid"] + ".png";
            }

        }
    }

    curl_easy_cleanup(c);
    return !image_url.empty();
}

// private

double SpeedTest::execute(const ServerInfo &server, const TestConfig &config, const opFn &pfunc, std::function<void(bool)> cb) {
    std::vector<std::thread> workers;
    double overall_speed = 0;
    std::mutex mtx;
    for (int i = 0; i < config.concurrency; i++) {
        workers.push_back(std::thread([&server, &overall_speed, &pfunc, &config, &mtx, cb](){
            long start_size = config.start_size;
            long max_size   = config.max_size;
            long incr_size  = config.incr_size;
            long curr_size  = start_size;

            auto spClient = SpeedTestClient(server);

            if (spClient.connect()) {
                long total_size = 0;
                long total_time = 0;
                auto start = std::chrono::steady_clock::now();
                std::vector<double> partial_results;
                while (curr_size < max_size){
                    long op_time = 0;
                    if ((spClient.*pfunc)(curr_size, config.buff_size, op_time)) {
                        total_size += curr_size;
                        total_time += op_time;
                        double metric = (curr_size * 8) / (static_cast<double>(op_time) / 1000);
                        partial_results.push_back(metric);
                        if (cb)
                            cb(true);
                    } else {
                        if (cb)
                            cb(false);
                    }
                    curr_size += incr_size;
                    auto stop = std::chrono::steady_clock::now();
                    if (std::chrono::duration_cast<std::chrono::milliseconds>(stop - start).count() > config.min_test_time_ms)
                        break;
                }

                spClient.close();
                std::sort(partial_results.begin(), partial_results.end());

                size_t skip = 0;
                size_t drop = 0;
                if (partial_results.size() >= 10){
                    skip = partial_results.size() / 4;
                    drop = 2;
                }

                size_t iter = 0;
                double real_sum = 0;
                for (auto it = partial_results.begin() + skip; it != partial_results.end() - drop; ++it ){
                    iter++;
                    real_sum += (*it);
                }
                mtx.lock();
                overall_speed += (real_sum / iter);
                mtx.unlock();
            } else {
                if (cb)
                    cb(false);
            }
        }));

    }
    for (auto &t : workers){
        t.join();
    }

    workers.clear();

    return overall_speed / 1000 / 1000;
}

template<typename T>
T SpeedTest::deg2rad(T n) {
    return (n * M_PI / 180);
}

template<typename T>
T SpeedTest::harversine(std::pair<T, T> n1, std::pair<T, T> n2) {
    T lat1r = deg2rad(n1.first);
    T lon1r = deg2rad(n1.second);
    T lat2r = deg2rad(n2.first);
    T lon2r = deg2rad(n2.second);
    T u = std::sin((lat2r - lat1r) / 2);
    T v = std::sin((lon2r - lon1r) / 2);
    return 2.0 * EARTH_RADIUS_KM * std::asin(std::sqrt(u * u + std::cos(lat1r) * std::cos(lat2r) * v * v));
}

CURLcode SpeedTest::httpGet(const std::string &url, std::stringstream &ss, CURL *handler, long timeout) {

    CURLcode code(CURLE_FAILED_INIT);
    CURL* curl = SpeedTest::curl_setup(handler);


    if (curl){
        if (CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_FILE, &ss))
            && CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_TIMEOUT, timeout))
            && CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_URL, url.c_str()))) {
            code = curl_easy_perform(curl);
        }
        if (handler == nullptr)
            curl_easy_cleanup(curl);
    }
    return code;
}

CURLcode SpeedTest::httpPost(const std::string &url, const std::string &postdata, std::stringstream &os, void *handler, long timeout) {

    CURLcode code(CURLE_FAILED_INIT);
    CURL* curl = SpeedTest::curl_setup(handler);

    if (curl){
        if (CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_FILE, &os))
            && CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_TIMEOUT, timeout))
            && CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_URL, url.c_str()))
            && CURLE_OK == (code = curl_easy_setopt(curl, CURLOPT_POSTFIELDS, postdata.c_str()))) {
            code = curl_easy_perform(curl);
        }
        if (handler == nullptr)
            curl_easy_cleanup(curl);
    }
    return code;
}

CURL *SpeedTest::curl_setup(CURL *handler) {
    CURL* curl = handler == nullptr ? curl_easy_init() : handler;
    if (curl){
        if (curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, &writeFunc) == CURLE_OK
            && curl_easy_setopt(curl, CURLOPT_NOPROGRESS, 1L) == CURLE_OK
            && curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L) == CURLE_OK
            && curl_easy_setopt(curl, CURLOPT_USERAGENT, SPEED_TEST_USER_AGENT) == CURLE_OK){
            return curl;
        } else {
            curl_easy_cleanup(handler);
            return nullptr;
        }
    }
    return nullptr;


}

size_t SpeedTest::writeFunc(void *buf, size_t size, size_t nmemb, void *userp) {

    if (userp){
        std::stringstream &os = *static_cast<std::stringstream *>(userp);
        std::streamsize len = size * nmemb;
        if(os.write(static_cast<char*>(buf), len))
            return  static_cast<size_t>(len);
    }
    return 0;
}

std::map<std::string, std::string> SpeedTest::parseQueryString(const std::string &query) {
    auto map = std::map<std::string, std::string>();
    auto pairs = splitString(query, '&');
    for (auto &p : pairs){
        auto kv = splitString(p, '=');
        if (kv.size() == 2){
            map[kv[0]] = kv[1];
        }
    }
    return map;
}

std::vector<std::string> SpeedTest::splitString(const std::string &instr, const char separator) {
    if (instr.empty())
        return std::vector<std::string>();

    std::vector<std::string> tokens;
    std::size_t start = 0, end = 0;
    while ((end = instr.find(separator, start)) != std::string::npos) {
        std::string temp = instr.substr(start, end - start);
        if (!temp.empty())
            tokens.push_back(temp);
        start = end + 1;
    }
    std::string temp = instr.substr(start);
    if (!temp.empty())
        tokens.push_back(temp);
    return tokens;

}

ServerInfo SpeedTest::processServerXMLNode(xmlTextReaderPtr reader) {

    auto name = xmlTextReaderConstName(reader);
    auto nodeName = std::string((char*)name);

    if (!name || nodeName != "server"){
        return ServerInfo();
    }

    if (xmlTextReaderAttributeCount(reader) > 0){
        auto info = ServerInfo();
        auto server_url         = xmlTextReaderGetAttribute(reader, BAD_CAST "url");
        auto server_lat         = xmlTextReaderGetAttribute(reader, BAD_CAST "lat");
        auto server_lon         = xmlTextReaderGetAttribute(reader, BAD_CAST "lon");
        auto server_name        = xmlTextReaderGetAttribute(reader, BAD_CAST "name");
        auto server_county      = xmlTextReaderGetAttribute(reader, BAD_CAST "country");
        auto server_cc          = xmlTextReaderGetAttribute(reader, BAD_CAST "cc");
        auto server_host        = xmlTextReaderGetAttribute(reader, BAD_CAST "host");
        auto server_id          = xmlTextReaderGetAttribute(reader, BAD_CAST "id");
        auto server_sponsor     = xmlTextReaderGetAttribute(reader, BAD_CAST "sponsor");

        if (server_name)
            info.name.append((char*)server_name);

        if (server_url)
            info.url.append((char*)server_url);

        if (server_county)
            info.country.append((char*)server_county);

        if (server_cc)
            info.country_code.append((char*)server_cc);

        if (server_host)
            info.host.append((char*)server_host);

        if (server_sponsor)
            info.sponsor.append((char*)server_sponsor);

        if (server_id)
            info.id  = std::atoi((char*)server_id);

        if (server_lat)
            info.lat = std::stof((char*)server_lat);

        if (server_lon)
            info.lon = std::stof((char*)server_lon);

        xmlFree(server_url);
        xmlFree(server_lat);
        xmlFree(server_lon);
        xmlFree(server_name);
        xmlFree(server_county);
        xmlFree(server_cc);
        xmlFree(server_host);
        xmlFree(server_id);
        xmlFree(server_sponsor);
        return info;
    }

    return ServerInfo();
}

bool SpeedTest::fetchServers(const std::string& url, std::vector<ServerInfo>& target, int &http_code) {
    std::stringstream oss;
    target.clear();

    auto isHttpSchema = url.find_first_of("http") == 0;

    CURL* curl = curl_easy_init();
    auto cres = httpGet(url, oss, curl, 20);

    if (cres != CURLE_OK)
        return false;

    if (isHttpSchema) {
        int req_status;
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &req_status);
        http_code = req_status;

        if (http_code != 200){
            curl_easy_cleanup(curl);
            return false;
        }
    } else {
        http_code = 200;
    }

    size_t len = oss.str().length();
    auto *xmlbuff = (char*)calloc(len + 1, sizeof(char));
    if (!xmlbuff){
        std::cerr << "Unable to calloc" << std::endl;
        curl_easy_cleanup(curl);
        return false;
    }
    memcpy(xmlbuff, oss.str().c_str(), len);
    oss.str("");

    xmlTextReaderPtr reader = xmlReaderForMemory(xmlbuff, static_cast<int>(len), nullptr, nullptr, 0);

    if (reader != nullptr) {
        IPInfo ipInfo;
        if (!SpeedTest::ipInfo(ipInfo)){
            curl_easy_cleanup(curl);
            free(xmlbuff);
            xmlFreeTextReader(reader);
            std::cerr << "OOPS!" <<std::endl;
            return false;
        }
        auto ret = xmlTextReaderRead(reader);
        while (ret == 1) {
            ServerInfo info = processServerXMLNode(reader);
            if (!info.url.empty()){
                info.distance = harversine(std::make_pair(ipInfo.lat, ipInfo.lon), std::make_pair(info.lat, info.lon));
                target.push_back(info);
            }
            ret = xmlTextReaderRead(reader);
        }
        xmlFreeTextReader(reader);
        if (ret != 0) {
            curl_easy_cleanup(curl);
            free(xmlbuff);
            std::cerr << "Failed to parse" << std::endl;
            return false;
        }
    } else {
        std::cerr << "Unable to initialize xml parser" << std::endl;
        curl_easy_cleanup(curl);
        free(xmlbuff);
        return false;
    }

    curl_easy_cleanup(curl);
    free(xmlbuff);
    xmlCleanupParser();
    std::sort(target.begin(), target.end(), [](const ServerInfo &a, const ServerInfo &b) -> bool {
        return a.distance < b.distance;
    });
    return true;
}

const ServerInfo SpeedTest::findBestServerWithin(const std::vector<ServerInfo> &serverList, long &latency,
                                                 const int sample_size, std::function<void(bool)> cb) {
    int i = sample_size;
    ServerInfo bestServer = serverList[0];

    latency = INT_MAX;

    for (auto &server : serverList){
        auto client = SpeedTestClient(server);

        if (!client.connect()){
            if (cb)
                cb(false);
            continue;
        }

        if (client.version() < mMinSupportedServer){
            client.close();
            continue;
        }

        long current_latency = LONG_MAX;
        if (testLatency(client, 20, current_latency)){
            if (current_latency < latency){
                latency = current_latency;
                bestServer = server;
            }
        }
        client.close();
        if (cb)
            cb(true);

        if (i-- < 0){
            break;
        }

    }
    return bestServer;
}

bool SpeedTest::testLatency(SpeedTestClient &client, const int sample_size, long &latency) {
    if (!client.connect()){
        return false;
    }
    latency = INT_MAX;
    long temp_latency = 0;
    for (int i = 0; i < sample_size; i++){
        if (client.ping(temp_latency)){
            if (temp_latency < latency){
                latency = temp_latency;
            }
        } else {
            return false;
        }
    }
    return true;
}

