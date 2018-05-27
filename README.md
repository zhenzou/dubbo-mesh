## dubbo-mesh

基于dubbo的Service Mesh Demo
2018阿里云中间件性能挑战赛参考

## 使用
```bash
docker network create --driver=bridge --subnet=10.10.10.0/24 --gateway=10.10.10.1 -o "com.docker.network.bridge.name"="benchmarker" -o "com.docker.network.bridge.enable_icc"="true" benchmarker

git clone git@github.com:zhenzou/dubbo-mesh.git $GOPATH/src/

cd $GOPATH/src/dubbo-mesh

docker build -t dubbo-mesh .

bash build/start-docker-mesh.sh

bash build/benchmark.sh # 需要clone官方benchmarker项目
```

## 说明
在第一天的提交就已经到了4K QPS了，后面的几天都是在调参，但并没有什么很大的突破。
目前来看最好的算法应该是WRR和加权最小连接数。
准备尝试一下加权RTT。

## 成绩
![成绩](doc/result.png)