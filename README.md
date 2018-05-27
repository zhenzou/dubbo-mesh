## dubbo-mesh

基于dubbo的Service Mesh Demo
2018阿里云中间件性能挑战赛参考

## 初始化
```bash
docker network create --driver=bridge --subnet=10.10.10.0/24 --gateway=10.10.10.1 -o "com.docker.network.bridge.name"="benchmarker" -o "com.docker.network.bridge.enable_icc"="true" benchmarker

docker build -t dubbo-mesh .
```
## 成绩
![成绩](doc/result.png)