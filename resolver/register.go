package resolver

import (
	"context"
	"path"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type ServiceInfo struct {
	Name     string
	Address  string
	Tags     []string
	Interval time.Duration
}

func NewServiceInfo() *ServiceInfo {
	return &ServiceInfo{
		Tags:     []string{"grpc"},
		Interval: 10,
	}
}

type Registry interface {
	Register(serviceInfo *ServiceInfo) (err error)
	Close() (err error)
}

type EtcdRegistry struct {
	etcdClient *clientv3.Client
	domain     string
	leaseID    clientv3.LeaseID
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewEtcdRegistry(ctx context.Context, etcdClient *clientv3.Client, domain string) *EtcdRegistry {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &EtcdRegistry{
		etcdClient: etcdClient,
		domain:     domain,
		ctx:        ctx,
		cancel:     cancelFunc,
	}
}

func (e *EtcdRegistry) Register(serviceInfo *ServiceInfo) (err error) {
	// 创建租约
	var lease *clientv3.LeaseGrantResponse
	lease, err = e.etcdClient.Grant(e.ctx, int64(serviceInfo.Interval))
	if err != nil {
		return
	}
	e.leaseID = lease.ID
	// 绑定租约
	target := path.Join(e.domain, serviceInfo.Name)
	var em endpoints.Manager
	em, err = endpoints.NewManager(e.etcdClient, target)
	if err != nil {
		return
	}
	key := path.Join(target, serviceInfo.Address)
	endpoint := endpoints.Endpoint{
		Addr: serviceInfo.Address,
		Metadata: map[string]string{
			"name": serviceInfo.Name,
			"tags": strings.Join(serviceInfo.Tags, ","),
		},
	}
	err = em.AddEndpoint(e.ctx, key, endpoint, clientv3.WithLease(e.leaseID))
	if err != nil {
		return
	}
	// 续租 发送心跳，表明服务正常
	var keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	keepAliveChan, err = e.etcdClient.KeepAlive(e.ctx, e.leaseID)
	if err != nil {
		return
	}
	// 监听续约
	go e.watcher(keepAliveChan)
	return
}

func (e *EtcdRegistry) watcher(resChan <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case l := <-resChan:
			if l == nil {
				return
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *EtcdRegistry) Close() (err error) {
	e.cancel()
	// 撤销租约
	e.etcdClient.Revoke(e.ctx, e.leaseID)
	return
}
