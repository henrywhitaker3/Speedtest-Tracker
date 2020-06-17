FROM ubuntu:latest

RUN mkdir -p /tmp/build /tmp/src
COPY *.h *.cpp *.h.in CMakeLists.txt /tmp/src/


RUN apt-get update && apt-get install -y g++ cmake make libcurl4-openssl-dev libxml2-dev libssl-dev && \
    cd /tmp/build && cmake -DCMAKE_BUILD_TYPE=Release ../src && make install && \
    apt-get remove --purge -y gcc make cmake libcurl4-openssl-dev libxml2-dev libssl-dev && \
    apt-get remove --purge -y `apt-mark showauto` && \
    apt-get install -y libxml2 libcurl3 && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/build /tmp/src

ENTRYPOINT ["/usr/local/bin/SpeedTest"]
