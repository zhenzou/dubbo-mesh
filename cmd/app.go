package cmd

import (
	"flag"
	"runtime"

	"dubbo-mesh/log"
	"dubbo-mesh/util"
)

var (
	CpuNum    int
	Port      int
	DubboPort int
	Etcd      string
	Memory    int
	Name      string
	Service   string
)

func init() {
	flag.IntVar(&CpuNum, "cpunum", runtime.NumCPU(), "指定cpu数量，默认CPU核数")
	flag.IntVar(&Memory, "m", 2048, "内存数，单位是M")
	flag.IntVar(&Port, "p", 20000, "监听端口，默认是20000,consumer端口")
	flag.IntVar(&DubboPort, "dp", 20880, "dubbo服务端口")
	flag.StringVar(&Etcd, "e", "http://127.0.0.1:2379", "Etcd 服务地址")
	flag.StringVar(&Name, "n", "consumer", "服务名称")
	flag.StringVar(&Service, "s", "com.alibaba.dubbo.performance.demo.provider.IHelloService", "dubbo服务全限定名")
}

func ParseFlag() {
	if !flag.Parsed() {
		flag.Parse()
	}
	runtime.GOMAXPROCS(CpuNum)
}

func Run(initFunc, jobFunc, cleanupFunc func() error) {
	ParseFlag()

	log.Infof("初始化 [%s]", Name)
	if err := initFunc(); err != nil {
		log.Infof("初始化 [%s] 失败：[%s]", Name, err)
		log.Panic(err)
	}
	log.Infof("初始化 [%s] 完成", Name)
	go func() {
		if err := jobFunc(); err != nil {
			log.Infof("[%s] 运行出错：[%v]", Name, err)
			log.Panic(err)
		}
	}()

	util.WaitForExitSign()
	log.Infof("[%s] 监听到退出信息，开始清理", Name)
	if err := cleanupFunc(); err != nil {
		log.Infof("[%s] 清理失败：[%v]", Name, err)
		log.Panic(err)
	}
	log.Infof("[%s] 清理完成，成功退出", Name)
	log.Sync()
}
