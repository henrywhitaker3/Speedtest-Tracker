//
// Created by Francesco Laurita on 6/3/16.
//

#ifndef SPEEDTEST_MD5UTIL_H
#define SPEEDTEST_MD5UTIL_H
#if defined(__APPLE__)
#  define COMMON_DIGEST_FOR_OPENSSL
#  include <CommonCrypto/CommonDigest.h>
#include <string>

#  define SHA1 CC_SHA1
#else
#  include <openssl/md5.h>
#endif

class MD5Util {
public:
    static std::string hexDigest(const std::string &str);
};


#endif //SPEEDTEST_MD5UTIL_H
