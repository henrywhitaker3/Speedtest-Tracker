//
// Created by Francesco Laurita on 6/8/16.
//

#ifndef SPEEDTEST_DATATYPES_H
#define SPEEDTEST_DATATYPES_H
#include <iostream>
#include <stdio.h>
#include <stdlib.h>
static const float EARTH_RADIUS_KM = 6371.0;

typedef struct ip_info_t {
    std::string ip_address;
    std::string isp;
    float lat;
    float lon;
} IPInfo;

typedef struct server_info_t {
    std::string url;
    std::string name;
    std::string country;
    std::string country_code;
    std::string host;
    std::string sponsor;
    int   id;
    float lat;
    float lon;
    float distance;

} ServerInfo;

typedef struct test_config_t {
    long start_size;
    long max_size;
    long incr_size;
    long buff_size;
    long min_test_time_ms;
    int  concurrency;
    std::string label;
} TestConfig;

#endif //SPEEDTEST_DATATYPES_H
