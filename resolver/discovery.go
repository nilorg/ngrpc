package resolver

import (
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	gresolver "google.golang.org/grpc/resolver"
)

type Discovery interface {
	Discover(domain string) (builder gresolver.Builder, err error)
	Close() (err error)
	Address() (addr string)
}

type EtcdDiscovery struct {
	etcdClient *clientv3.Client
	builder    gresolver.Builder
	domain     string
	service    string
}

func NewEtcdDiscovery(etcdClient *clientv3.Client, domain string) *EtcdDiscovery {
	builder, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	return &EtcdDiscovery{
		etcdClient: etcdClient,
		builder:    builder,
		domain:     domain,
	}
}

func (e *EtcdDiscovery) Discover(service string) (builder gresolver.Builder, err error) {
	e.service = service
	builder = e.builder
	return
}

func (e *EtcdDiscovery) Close() (err error) {
	return
}

func (e *EtcdDiscovery) Address() (addr string) {
	addr = fmt.Sprintf("%s:///%s/%s", e.builder.Scheme(), e.domain, e.service)
	return
}
