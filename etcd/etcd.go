package etcd

import (
	"context"
	"time"

	etcdV3 "go.etcd.io/etcd/client/v3"
)

type EtcdHelper struct {
	etcd *etcdV3.Client
}

func NewEtcdHelper(clientEndpoints []string) (*EtcdHelper, error) {
	client, err := etcdV3.New(etcdV3.Config{
		Endpoints:   clientEndpoints,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdHelper{
		etcd: client,
	}, nil
}

func (et *EtcdHelper) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	rep, err := et.etcd.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}
	rep.OpResponse()
	return nil
}

func (et *EtcdHelper) Close(){
	et.etcd.Close()
}