//
// Created by Francesco Laurita on 5/29/16.
//

#ifndef SPEEDTEST_SPEEDTEST_H
#define SPEEDTEST_SPEEDTEST_H

#include "SpeedTestConfig.h"
#include "SpeedTestClient.h"
#include <libxml/xmlreader.h>
#include <functional>
#include <cmath>
#include <curl/curl.h>
#include <fstream>
#include <sstream>
#include <iostream>
#include <map>
#include <vector>
#include <algorithm>
#include <thread>
#include <mutex>
#include "DataTypes.h"

class SpeedTestClient;
typedef bool (SpeedTestClient::*opFn)(const long size, const long chunk_size, long &millisec);
typedef void (*progressFn)(bool success);


class SpeedTest {
public:
    explicit SpeedTest(float minServerVersion);
    ~SpeedTest();
    CURLcode httpGet(const std::string& url, std::stringstream& os, CURL *handler = nullptr, long timeout = 30);
    CURLcode httpPost(const std::string& url, const std::string& postdata, std::stringstream& os, CURL *handler = nullptr, long timeout = 30);
    static std::map<std::string, std::string> parseQueryString(const std::string& query);
    static std::vector<std::string> splitString(const std::string& instr, char separator);
    bool ipInfo(IPInfo& info);
    const std::vector<ServerInfo>& serverList();
    const ServerInfo bestServer(int sample_size = 5, std::function<void(bool)> cb = nullptr);
    bool setServer(ServerInfo& server);
    const long &latency();
    bool downloadSpeed(const ServerInfo& server, const TestConfig& config, double& result, std::function<void(bool)> cb = nullptr);
    bool uploadSpeed(const ServerInfo& server, const TestConfig& config, double& result, std::function<void(bool)> cb = nullptr);
    bool jitter(const ServerInfo& server, long& result, int sample = 40);
    bool share(const ServerInfo& server, std::string& image_url);
private:
    bool fetchServers(const std::string& url,  std::vector<ServerInfo>& target, int &http_code);
    bool testLatency(SpeedTestClient& client, int sample_size, long& latency);
    const ServerInfo findBestServerWithin(const std::vector<ServerInfo>& serverList, long& latency, int sample_size = 5, std::function<void(bool)> cb = nullptr);
    static CURL* curl_setup(CURL* curl = nullptr);
    static size_t writeFunc(void* buf, size_t size, size_t nmemb, void* userp);
    static ServerInfo processServerXMLNode(xmlTextReaderPtr reader);
    double execute(const ServerInfo &server, const TestConfig &config, const opFn &fnc, std::function<void(bool)> cb = nullptr);
    template <typename T>
        static T deg2rad(T n);
    template <typename T>
        static T harversine(std::pair<T, T> n1, std::pair<T, T> n2);

    IPInfo mIpInfo;
    std::vector<ServerInfo> mServerList;
    long mLatency;
    double mUploadSpeed;
    double mDownloadSpeed;
    float mMinSupportedServer;

};


#endif //SPEEDTEST_SPEEDTEST_H
