package registry

import (
	"fmt"
	"context"
	"strings"
	"strconv"

	etcd "github.com/coreos/etcd/clientv3"

	"dubbo-mesh/util"
	"dubbo-mesh/log"
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
	resp, err := client.Grant(context.Background(), 30)
	if err != nil {
		panic(err)
	}
	leaseID := resp.ID
	etcd := &Etcd{leaseID, client}
	go etcd.keepalive()
	return etcd
}

type Etcd struct {
	leaseId etcd.LeaseID
	client  *etcd.Client
}

func (this *Etcd) keepalive() error {
	ch, err := this.client.Lease.KeepAlive(context.Background(), this.leaseId)
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
	_, err := this.client.Put(context.Background(), key, "", etcd.WithLease(this.leaseId))
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
		log.Debug("get service :", key)
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
		endpoint := &Endpoint{Host: split[0], Port: port}
		log.Debug("endpoint:", endpoint.String())
		endpoints = append(endpoints, endpoint)
	}
	return
}
