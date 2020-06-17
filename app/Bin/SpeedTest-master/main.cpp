#include <iostream>
#include <map>
#include <iomanip>
#include "SpeedTest.h"
#include "TestConfigTemplate.h"
#include "CmdOptions.h"
#include <csignal>

void banner(){
    std::cout << "SpeedTest++ version " << SpeedTest_VERSION_MAJOR << "." << SpeedTest_VERSION_MINOR << std::endl;
    std::cout << "Speedtest.net command line interface" << std::endl;
    std::cout << "Info: " << SpeedTest_HOME_PAGE << std::endl;
    std::cout << "Author: " << SpeedTest_AUTHOR << std::endl;
}

void usage(const char* name){
    std::cerr << "Usage: " << name << " ";
    std::cerr << " [--latency] [--quality] [--download] [--upload] [--share] [--help]\n"
            "      [--test-server host:port] [--output verbose|text|json]\n";
    std::cerr << "optional arguments:" << std::endl;
    std::cerr << "  --help                      Show this message and exit\n";
    std::cerr << "  --latency                   Perform latency test only\n";
    std::cerr << "  --download                  Perform download test only. It includes latency test\n";
    std::cerr << "  --upload                    Perform upload test only. It includes latency test\n";
    std::cerr << "  --share                     Generate and provide a URL to the speedtest.net share results image\n";
    std::cerr << "  --test-server host:port     Run speed test against a specific server\n";
    std::cerr << "  --output verbose|text|json  Set output type. Default: verbose\n";
}

int main(const int argc, const char **argv) {

    ProgramOptions programOptions;

    if (!ParseOptions(argc, argv, programOptions)){
        usage(argv[0]);
        return EXIT_FAILURE;
    }

    if (programOptions.output_type == OutputType::verbose){
        banner();
        std::cout << std::endl;
    }


    if (programOptions.help) {
        usage(argv[0]);
        return EXIT_SUCCESS;
    }

    signal(SIGPIPE, SIG_IGN);


    auto sp = SpeedTest(SPEED_TEST_MIN_SERVER_VERSION);
    IPInfo info;
    ServerInfo serverInfo;
    ServerInfo serverQualityInfo;

    if (programOptions.output_type == OutputType::json)
        std::cout << "{";

    if (!sp.ipInfo(info)){
        std::cerr << "Unable to retrieve your IP info. Try again later" << std::endl;
        if (programOptions.output_type == OutputType::json)
            std::cout << "\"error\":\"unable to retrieve your ip info\"}" << std::endl;
        return EXIT_FAILURE;
    }

    if (programOptions.output_type == OutputType::verbose){
        std::cout << "IP: " << info.ip_address
                  << " ( " << info.isp << " ) "
                  << "Location: [" << info.lat << ", " << info.lon << "]" << std::endl;
    } else if (programOptions.output_type == OutputType::text) {
        std::cout << "IP=" << info.ip_address << std::endl;
        std::cout << "IP_LAT=" << info.lat << std::endl;
        std::cout << "IP_LON=" << info.lon << std::endl;
        std::cout << "PROVIDER=" << info.isp << std::endl;
    } else if (programOptions.output_type == OutputType::json) {
        std::cout << "\"client\":{";
        std::cout << "\"ip\":\""  << info.ip_address << "\",";
        std::cout << "\"lat\":\"" << info.lat << "\",";
        std::cout << "\"lon\":\"" << info.lon << "\",";
        std::cout << "\"isp\":\"" << info.isp << "\"";
        std::cout << "},";
    }

    auto serverList = sp.serverList();

    if (programOptions.selected_server.empty()){
        if (programOptions.output_type == OutputType::verbose)
            std::cout << "Finding fastest server... " << std::flush;

        if (serverList.empty()){
            std::cerr << "Unable to download server list. Try again later" << std::endl;
            if (programOptions.output_type == OutputType::json)
                std::cout << "\"error\":\"unable to download server list\"}" << std::endl;
            return EXIT_FAILURE;
        }

        if (programOptions.output_type == OutputType::verbose)
            std::cout << serverList.size() << " Servers online" << std::endl;
        else if (programOptions.output_type == OutputType::json)
            std::cout << "\"servers_online\":\"" << serverList.size() << "\",";


        serverInfo = sp.bestServer(10, [&programOptions](bool success) {
            if (programOptions.output_type == OutputType::verbose)
                std::cout << (success ? '.' : '*') << std::flush;
        });

        if (programOptions.output_type == OutputType::verbose){
            std::cout << std::endl;
            std::cout << "Server: " << serverInfo.name
                      << " " << serverInfo.host
                      << " by " << serverInfo.sponsor
                      << " (" << serverInfo.distance << " km from you): "
                      << sp.latency() << " ms" << std::endl;
        } else if (programOptions.output_type == OutputType::text) {
            std::cout << "TEST_SERVER_HOST=" << serverInfo.host << std::endl;
            std::cout << "TEST_SERVER_DISTANCE=" << serverInfo.distance << std::endl;

        }
        else if (programOptions.output_type == OutputType::json) {
            std::cout << "\"server\":{";
            std::cout << "\"name\":\"" << serverInfo.name << "\",";
            std::cout << "\"sponsor\":\"" << serverInfo.sponsor << "\",";
            std::cout << "\"distance\":\"" << serverInfo.distance << "\",";
            std::cout << "\"latency\":\"" << sp.latency() << "\",";
            std::cout << "\"host\":\"" << serverInfo.host << "\"";
            std::cout << "},";
        }

    } else {

        serverInfo.host.append(programOptions.selected_server);
        sp.setServer(serverInfo);

        for (auto &s : serverList) {
            if (s.host == serverInfo.host)
                serverInfo.id = s.id;
        }

        if (programOptions.output_type == OutputType::verbose)
            std::cout << "Selected server: " << serverInfo.host << std::endl;
        else if (programOptions.output_type == OutputType::text) {
            std::cout << "TEST_SERVER_HOST=" << serverInfo.host << std::endl;
        }
        else if (programOptions.output_type == OutputType::json) {
            std::cout << "\"server\":{";
            std::cout << "\"host\":\"" << serverInfo.host << "\"";
            std::cout << "},";
        }
    }

    if (programOptions.output_type == OutputType::verbose)
        std::cout << "Ping: " << sp.latency() << " ms." << std::endl;
    else if (programOptions.output_type == OutputType::text)
        std::cout << "LATENCY=" << sp.latency() << std::endl;
    else if (programOptions.output_type == OutputType::json) {
        std::cout << "\"ping\":\"";
        std::cout << std::fixed;
        std::cout << sp.latency() << "\",";
    }

    long jitter = 0;
    if (programOptions.output_type == OutputType::verbose)
        std::cout << "Jitter: " << std::flush;
    if (sp.jitter(serverInfo, jitter)){
        if (programOptions.output_type == OutputType::verbose)
            std::cout << jitter << " ms." << std::endl;
        else if (programOptions.output_type == OutputType::text)
            std::cout << "JITTER=" << jitter << std::endl;
        else if (programOptions.output_type == OutputType::json) {
            std::cout << "\"jitter\":\"";
            std::cout << std::fixed;
            std::cout << jitter << "\",";
        }
    } else {
        std::cerr << "Jitter measurement is unavailable at this time." << std::endl;
    }

    if (programOptions.latency) {
        if (programOptions.output_type == OutputType::json)
            std::cout << "\"_\":\"only latency requested\"}" << std::endl;
        return EXIT_SUCCESS;
    }


    if (programOptions.output_type == OutputType::verbose)
        std::cout << "Determine line type (" << preflightConfigDownload.concurrency << ") "  << std::flush;
    double preSpeed = 0;
    if (!sp.downloadSpeed(serverInfo, preflightConfigDownload, preSpeed, [&programOptions](bool success){
        if (programOptions.output_type == OutputType::verbose)
            std::cout << (success ? '.' : '*') << std::flush;
    })){
        std::cerr << "Pre-flight check failed." << std::endl;
        if (programOptions.output_type == OutputType::json)
            std::cout << "\"error\":\"pre-flight check failed\"}" << std::endl;
        return EXIT_FAILURE;
    }

    if (programOptions.output_type == OutputType::verbose)
        std::cout << std::endl;

    TestConfig uploadConfig;
    TestConfig downloadConfig;
    testConfigSelector(preSpeed, uploadConfig, downloadConfig);

    if (programOptions.output_type == OutputType::verbose)
        std::cout << downloadConfig.label << std::endl;


    if (!programOptions.upload){
        if (programOptions.output_type == OutputType::verbose){
            std::cout << std::endl;
            std::cout << "Testing download speed (" << downloadConfig.concurrency << ") "  << std::flush;
        }

        double downloadSpeed = 0;
        if (sp.downloadSpeed(serverInfo, downloadConfig, downloadSpeed, [&programOptions](bool success){
            if (programOptions.output_type == OutputType::verbose)
                std::cout << (success ? '.' : '*') << std::flush;
        })){
            if (programOptions.output_type == OutputType::verbose){
                std::cout << std::endl;
                std::cout << "Download: ";
                std::cout << std::fixed;
                std::cout << std::setprecision(2);
                std::cout << downloadSpeed << " Mbit/s" << std::endl;
            } else if (programOptions.output_type == OutputType::text) {
                std::cout << "DOWNLOAD_SPEED=";
                std::cout << std::fixed;
                std::cout << std::setprecision(2);
                std::cout << downloadSpeed << std::endl;
            } else if (programOptions.output_type == OutputType::json) {
                std::cout << "\"download\":\"";
                std::cout << std::fixed;
                std::cout << (downloadSpeed*1000*1000) << "\",";
            }
        } else {
            std::cerr << "Download test failed." << std::endl;
            if (programOptions.output_type == OutputType::json)
                std::cout << "\"error\":\"download test failed\"}" << std::endl;
            return EXIT_FAILURE;
        }
    }

    if (programOptions.download) {
        if (programOptions.output_type == OutputType::json)
            std::cout << "\"_\":\"only download requested\"}" << std::endl;
        return EXIT_SUCCESS;
    }

    if (programOptions.output_type == OutputType::verbose)
        std::cout << "Testing upload speed (" << uploadConfig.concurrency << ") "  << std::flush;

    double uploadSpeed = 0;
    if (sp.uploadSpeed(serverInfo, uploadConfig, uploadSpeed, [&programOptions](bool success){
        if (programOptions.output_type == OutputType::verbose)
            std::cout << (success ? '.' : '*') << std::flush;
    })){
        if (programOptions.output_type == OutputType::verbose){
            std::cout << std::endl;
            std::cout << "Upload: ";
            std::cout << std::fixed;
            std::cout << std::setprecision(2);
            std::cout << uploadSpeed << " Mbit/s" << std::endl;
        } else if (programOptions.output_type == OutputType::text) {
            std::cout << "UPLOAD_SPEED=";
            std::cout << std::fixed;
            std::cout << std::setprecision(2);
            std::cout << uploadSpeed << std::endl;
        } else if (programOptions.output_type == OutputType::json) {
            std::cout << "\"upload\":\"";
            std::cout << std::fixed;
            std::cout << (uploadSpeed*1000*1000) << "\",";
        }

    } else {
        std::cerr << "Upload test failed." << std::endl;
        if (programOptions.output_type == OutputType::json)
            std::cout << "\"error\":\"upload test failed\"}" << std::endl;
        return EXIT_FAILURE;
    }


    if (programOptions.share){
        std::string share_it;
        if (sp.share(serverInfo, share_it)) {
            if (programOptions.output_type == OutputType::verbose) {
                std::cout << "Results image: " << share_it << std::endl;
            } else if (programOptions.output_type == OutputType::text) {
                std::cout << "IMAGE_URL=" << share_it << std::endl;
            } else if (programOptions.output_type == OutputType::json) {
                std::cout << "\"share\":\"" << share_it << "\",";
            }
        }
    }

    if (programOptions.output_type == OutputType::json)
        std::cout << "\"_\":\"all ok\"}" << std::endl;

    return EXIT_SUCCESS;
}
