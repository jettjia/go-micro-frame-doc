FROM golang:latest

ADD ./nacos-coredns-plugin /tmp/nacos-coredns-plugin
ADD ./coredns /tmp/coredns

# 编译构建
RUN sh /tmp/nacos-coredns-plugin/bin/build.sh 

# 导出默认端口，tcp和udp的都要导出
EXPOSE 53 53/udp

VOLUME ["/etc/coredns/Corefile"]

CMD ["/go/src/coredns/coredns", "-conf", "/etc/coredns/Corefile"]
