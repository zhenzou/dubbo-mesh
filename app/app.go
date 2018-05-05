package app

import (
	"flag"
	"runtime"
	"strings"

	"dubbo-mesh/log"
	"dubbo-mesh/util"
)

var (
	CpuNum int
	Name   string
)

func init() {
	flag.IntVar(&CpuNum, "cpunum", runtime.NumCPU(), "指定cpu数量，默认CPU核数")
}

func ParseFlag() {
	if !flag.Parsed() {
		flag.Parse()
	}
	runtime.GOMAXPROCS(CpuNum)
}

func Run(appName string, initFunc, jobFunc, cleanupFunc func() error) {
	ParseFlag()
	Name = appName
	log.Infof("running %s in %s mode", Name, Mode.String())

	log.Infof("初始化 [%s]", appName)
	if err := initFunc(); err != nil {
		log.Infof("初始化 [%s] 失败：[%s]", appName, err)
		panic(err)
	}
	log.Infof("初始化 [%s] 完成", appName)
	go func() {
		if err := jobFunc(); err != nil {
			log.Infof("[%s] 运行出错：[%v]", appName, err)
			panic(err)
		}
	}()

	util.WaitForExitSign()
	log.Infof("[%s] 监听到退出信息，开始清理", appName)
	if err := cleanupFunc(); err != nil {
		log.Infof("[%s] 清理失败：[%v]", appName, err)
		panic(err)
	}
	log.Infof("[%s] 清理完成，成功退出", appName)
	log.Sync()
}

func Funcs(funcs ...func() error) func() error {
	return func() error {
		for _, fun := range funcs {
			if err := fun(); err != nil {
				return err
			}
		}
		return nil
	}
}

func LogWrapper(msg string, fun func() error) func() error {
	return func() error {
		log.Info(msg + " 开始")
		if err := fun(); err != nil {
			log.Infof("%s 失败:%v", msg, err)
			return err
		}
		log.Info(msg + " 完成")
		return nil
	}
}

type RunningMode string

const (
	DevMode  RunningMode = "dev"
	ProdMode RunningMode = "prod"
	TestMode RunningMode = "test"
)

func (this RunningMode) String() string {
	return strings.ToLower(string(this))
}
