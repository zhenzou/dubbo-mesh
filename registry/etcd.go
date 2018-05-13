package registry

import (
	"fmt"
	"context"
	"strings"
	"strconv"
	"runtime"

	etcd "github.com/coreos/etcd/clientv3"

	"dubbo-mesh/util"
	"dubbo-mesh/log"
	"dubbo-mesh/json"
	"dubbo-mesh/cmd"
)

func NewEtcdFromAddr(addr string) Registry {
	endpoints := strings.Split(addr, ",")
	cfg := etcd.Config{
		Endpoints: endpoints,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		panic(err)
	}
	return NewEtcd(client)
}

func NewEtcd(client *etcd.Client) Registry {
	etcd := &Etcd{client: client}
	return etcd
}

type Etcd struct {
	client *etcd.Client
}

func (this *Etcd) keepalive(leaseId etcd.LeaseID) error {
	ch, err := this.client.Lease.KeepAlive(context.Background(), leaseId)
	if err != nil {
		return err
	}
	for resp := range ch {
		log.Info("keepalive ", resp.ID)
	}
	return nil
}

func (this *Etcd) Register(serviceName string, port int) error {
	key := this.strKey(serviceName, port)
	resp, err := this.client.Grant(context.Background(), 30)
	if err != nil {
		panic(err)
	}
	leaseId := resp.ID
	go this.keepalive(leaseId)
	_, err = this.client.Put(context.Background(), key, this.system(), etcd.WithLease(leaseId))
	if err != nil {
		return err
	}
	log.Info("register a new service:", key)
	return nil
}

func (this *Etcd) strKey(serviceName string, port int) string {
	key := fmt.Sprintf("/%s/%s/%s:%d", RootPath, serviceName, util.LocalIp(), port)
	return key
}

func (this *Etcd) system() string {
	mem := &runtime.MemStats{}
	system := &System{
		CpuNum:      runtime.NumCPU(),
		TotalMemory: cmd.Memory,
		UsedMemory:  int(mem.Mallocs),
	}
	bytes, _ := json.Marshal(system)
	return util.BytesToString(bytes)
}

func (this *Etcd) prefix(serviceName string) string {
	key := fmt.Sprintf("/%s/%s/", RootPath, serviceName)
	return key
}

func (this *Etcd) Find(serviceName string) (endpoints []*Endpoint, err error) {
	prefix := this.prefix(serviceName)
	resp, err := this.client.Get(context.Background(), prefix, etcd.WithPrefix())
	if err != nil {
		return
	}
	endpoints = make([]*Endpoint, 0, resp.Count)
	for _, kv := range resp.Kvs {
		key := string(kv.Key)
		addr := strings.TrimPrefix(key, prefix)
		split := strings.Split(addr, ":")
		if len(split) != 2 {
			log.Warn("get wrong service ", key)
			continue
		}
		port, err := strconv.Atoi(split[1])
		if err != nil {
			log.Warn("get wrong service ", key)
			continue
		}
		system := &System{}
		json.Unmarshal(kv.Value, system)
		endpoint := &Endpoint{Host: split[0], Port: port, System: system}
		endpoints = append(endpoints, endpoint)
	}
	return
}
