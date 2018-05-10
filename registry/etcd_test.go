package registry

import (
	"reflect"
	"testing"

	etcd "github.com/coreos/etcd/clientv3"
)

var (
	cfg = etcd.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	client   *etcd.Client
	registry Registry
)

func init() {
	var err error
	client, err = etcd.New(cfg)
	if err != nil {
		panic(err)
	}
	registry = NewEtcd(client)
}

func TestEtcd_Register(t *testing.T) {

	type args struct {
		serviceName string
		port        int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "service_hello",
			args:    args{"service_hello", 1024},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := registry.Register(tt.args.serviceName, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Etcd.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEtcd_Find(t *testing.T) {

	type args struct {
		serviceName string
	}
	tests := []struct {
		name          string
		args          args
		wantEndpoints []*Endpoint
		wantErr       bool
	}{
		{
			name:          "service_hello",
			args:          args{"service_hello"},
			wantEndpoints: []*Endpoint{&Endpoint{Host: "192.168.3.2", Port: 1024}},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotEndpoints, err := registry.Find(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Etcd.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEndpoints, tt.wantEndpoints) {
				t.Errorf("Etcd.Find() = %v, want %v", gotEndpoints, tt.wantEndpoints)
			}
		})
	}
}
