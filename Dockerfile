# Builder container
FROM registry.cn-hangzhou.aliyuncs.com/aliware2018/services AS builder
FROM golang

WORKDIR $GOPATH/src/dubbo-mesh
ADD . $GOPATH/src/dubbo-mesh
RUN set -ex && bash build/build_cmd.sh dev all

# Runner container

COPY --from=builder /root/workspace/services/mesh-provider/target/mesh-provider-1.0-SNAPSHOT.jar /root/dists/mesh-provider.jar
COPY --from=builder /root/workspace/services/mesh-consumer/target/mesh-consumer-1.0-SNAPSHOT.jar /root/dists/mesh-consumer.jar
COPY --from=builder /usr/local/bin/docker-entrypoint.sh /usr/local/bin

COPY build/start-agent.sh /usr/local/bin
COPY build/dist/consumer /root/dists/consumer
COPY build/dist/provider /root/dists/provider

RUN set -ex && chmod a+x /usr/local/bin/start-agent.sh && mkdir -p /root/logs

ENTRYPOINT ["docker-entrypoint.sh"]
