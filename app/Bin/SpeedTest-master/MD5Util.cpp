//
// Created by Francesco Laurita on 6/3/16.
//

#include <sstream>
#include "MD5Util.h"

std::string MD5Util::hexDigest(const std::string &str) {
    unsigned char digest[MD5_DIGEST_LENGTH];

    MD5_CTX ctx;
    MD5_Init(&ctx);
    MD5_Update(&ctx, str.c_str(), str.size());
    MD5_Final(digest, &ctx);

    char hexDigest[33] = {'\0'};
    for (int i = 0; i < 16; i++)
        std::sprintf(&hexDigest[i*2], "%02x", (unsigned int)digest[i]);

    return std::string(hexDigest);
}

