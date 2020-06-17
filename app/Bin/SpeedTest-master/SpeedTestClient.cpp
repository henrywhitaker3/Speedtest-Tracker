//
// Created by Francesco Laurita on 5/30/16.
//

#include <arpa/inet.h>
#include <netdb.h>
#include "SpeedTestClient.h"

SpeedTestClient::SpeedTestClient(const ServerInfo &serverInfo): mServerInfo(serverInfo),
                                                                                  mSocketFd(0),
                                                                                  mServerVersion(-1.0){}
SpeedTestClient::~SpeedTestClient() {
    close();
}

// It returns current timestamp in ms


// It connects and initiates client/server handshaking
bool SpeedTestClient::connect() {

    if (mSocketFd){
        return true;
    }

    auto ret = mkSocket();
    if (!ret)
        return ret;

    std::string reply;

    if (!SpeedTestClient::writeLine(mSocketFd, "HI")){
        close();
        return false;
    }


    if (SpeedTestClient::readLine(mSocketFd, reply)){
        std::stringstream reply_stream(reply);
        std::string hello;
        reply_stream >> hello >> mServerVersion;
        if (reply_stream.fail()) {
            close();
            return false;
        }

        if (!reply.empty() && "HELLO" == hello){
            return true;
        }

    }
    close();
    return false;
}

// It closes a connection
void SpeedTestClient::close() {
    if (mSocketFd){
        SpeedTestClient::writeLine(mSocketFd, "QUIT");
        ::close(mSocketFd);
    }

}

// It executes PING command
bool SpeedTestClient::ping(long &millisec) {
    if (!mSocketFd)
        return false;

    std::stringstream cmd;
    std::string reply;

    auto start = std::chrono::steady_clock::now();
    cmd << "PING " << start.time_since_epoch().count();

    if (!SpeedTestClient::writeLine(mSocketFd, cmd.str())){
        return false;
    }

    if (SpeedTestClient::readLine(mSocketFd, reply)){
        if (reply.substr(0, 5) == "PONG "){
            auto stop = std::chrono::steady_clock::now();
            millisec = std::chrono::duration_cast<std::chrono::milliseconds>(stop - start).count();
            return true;
        }
    }

    close();
    return false;
}

// It executes DOWNLOAD command
bool SpeedTestClient::download(const long size, const long chunk_size, long &millisec) {
    std::stringstream cmd;
    cmd << "DOWNLOAD " << size;

    if (!SpeedTestClient::writeLine(mSocketFd, cmd.str())){
        return false;
    }


    char *buff = new char[chunk_size];
    for (size_t i = 0; i < static_cast<size_t>(chunk_size); i++)
        buff[i] = '\0';

    long missing = 0;
    auto start = std::chrono::steady_clock::now();
    while (missing != size){
        auto current = read(mSocketFd, buff, static_cast<size_t>(chunk_size));

        if (current <= 0){
            delete[] buff;
            return false;
        }
        missing += current;
    }

    auto stop = std::chrono::steady_clock::now();
    millisec = std::chrono::duration_cast<std::chrono::milliseconds>(stop - start).count();
    delete[] buff;
    return true;
}

// It executes UPLOAD command
bool SpeedTestClient::upload(const long size, const long chunk_size, long &millisec) {
    std::stringstream cmd;
    cmd << "UPLOAD " << size << "\n";
    auto cmd_len = cmd.str().length();

    char *buff = new char[chunk_size];
    for(size_t i = 0; i < static_cast<size_t>(chunk_size); i++)
        buff[i] = static_cast<char>(rand() % 256);

    long missing = size;
    auto start = std::chrono::steady_clock::now();

    if (!SpeedTestClient::writeLine(mSocketFd, cmd.str())){
        delete[] buff;
        return false;
    }

    ssize_t w = cmd_len;
    missing -= w;

    while(missing > 0){
        if (missing - chunk_size > 0){
            w = write(mSocketFd, buff, static_cast<size_t>(chunk_size));
            if (w != chunk_size){
                delete[] buff;
                return false;
            }
            missing -= w;
        } else {
            buff[missing - 1] = '\n';
            w = write(mSocketFd, buff, static_cast<size_t>(missing));
            if (w != missing){
                delete[] buff;
                return false;
            }
            missing -= w;
        }

    }
    std::string reply;
    if (!SpeedTestClient::readLine(mSocketFd, reply)){
        delete[] buff;
        return false;
    }
    auto stop = std::chrono::steady_clock::now();

    std::stringstream ss;
    ss << "OK " << size << " ";
    millisec = std::chrono::duration_cast<std::chrono::milliseconds>(stop - start).count();
    delete[] buff;
    return reply.substr(0, ss.str().length()) == ss.str();

}

bool SpeedTestClient::mkSocket() {
    mSocketFd = socket(AF_INET, SOCK_STREAM, 0);

    if (!mSocketFd){
        return false;
    }

    auto hostp = hostport();
#if __APPLE__
    struct hostent *server = gethostbyname(hostp.first.c_str());
    if (server == nullptr) {
        return false;
    }
#else
    struct hostent server;
    char tmpbuf[BUFSIZ];
    struct hostent *result;
    int errnop;
    if (gethostbyname_r(hostp.first.c_str(), &server, (char *)&tmpbuf, BUFSIZ, &result, &errnop)) {
        return false;
    }
#endif

    int portno = hostp.second;
    struct sockaddr_in serv_addr{};
    memset(&serv_addr, 0, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;

#if __APPLE__
    memcpy(&serv_addr.sin_addr.s_addr, server->h_addr, (size_t)server->h_length);
#else
    memcpy(&serv_addr.sin_addr.s_addr, server.h_addr, (size_t)server.h_length);
#endif

    serv_addr.sin_port = htons(static_cast<uint16_t>(portno));

    /* Dial */
    return ::connect(mSocketFd, (struct sockaddr*)&serv_addr, sizeof(serv_addr)) >= 0;
}


float SpeedTestClient::version() {
    return mServerVersion;
}

const std::pair<std::string, int> SpeedTestClient::hostport() {
    std::string targetHost = mServerInfo.host;
    std::size_t found = targetHost.find(':');
    std::string host  = targetHost.substr(0, found);
    std::string port  = targetHost.substr(found + 1, targetHost.length() - found);
    return std::pair<std::string, int>(host, std::atoi(port.c_str()));
}

bool SpeedTestClient::readLine(int &fd, std::string &buffer) {
    buffer.clear();
    if (!fd)
        return false;
    char c;
    while(true){
        auto n = read(fd, &c, 1);
        if (n == -1)
            return false;
        if (c == '\n' || c == '\r')
            break;

        buffer += c;

    }
    return true;
}

bool SpeedTestClient::writeLine(int &fd, const std::string &buffer) {
    if (!fd)
        return false;

    auto len = static_cast<ssize_t >(buffer.length());
    if (len == 0)
        return false;

    std::string buff_copy = buffer;

    if (buff_copy.find_first_of('\n') == std::string::npos){
        buff_copy += '\n';
        len += 1;
    }
    auto n = write(fd, buff_copy.c_str(), len);
    return n == len;
}




