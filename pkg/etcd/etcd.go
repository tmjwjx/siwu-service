package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type EtcdClient struct {
	client     *clientv3.Client
	kv         clientv3.KV
	lease      clientv3.Lease
	watch      clientv3.Watcher
	serverList map[string]string // 存储服务名到地址的映射（服务名唯一）
	lock       sync.Mutex
}

// NewEtcdClient  初始化客户端对象
func NewEtcdClient(addr []string) (*EtcdClient, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		fmt.Printf("create connection etcd failed %s\n", err)
		return nil, err
	}

	// 得到 KV 、Lease、 Watcher 的API子集
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	watch := clientv3.NewWatcher(client)

	// 给客户端对象赋值
	c := &EtcdClient{
		client:     client,
		kv:         kv,
		lease:      lease,
		watch:      watch,
		serverList: make(map[string]string),
	}
	return c, nil
}

// 根据服务名获取单个实例地址
func (c *EtcdClient) getServiceAddress(serviceName string) (string, error) {
	// 直接读取键值，不使用 WithPrefix
	resp, err := c.kv.Get(context.Background(), serviceName)
	if err != nil {
		fmt.Printf("getServiceAddress failed: %s\n", err)
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("服务 %s 不存在", serviceName)
	}

	// 提取第一个键值对的地址（唯一）
	addr := string(resp.Kvs[0].Value)
	c.SetServiceList(serviceName, addr) // 更新本地缓存
	return addr, nil
}

// SetServiceList 设置服务地址（简化）
func (c *EtcdClient) SetServiceList(serviceName, addr string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.serverList[serviceName] = addr
	fmt.Printf("更新服务缓存: 服务名=%s, 地址=%s\n", serviceName, addr)
}

// DelServiceList 删除服务地址
func (c *EtcdClient) DelServiceList(serviceName string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.serverList, serviceName)
	fmt.Printf("删除服务缓存: 服务名=%s\n", serviceName)
}

// GetService 获取服务地址（对外接口）
func (c *EtcdClient) GetService(serviceName string) (string, error) {
	addr, err := c.getServiceAddress(serviceName)
	if err != nil {
		return "", fmt.Errorf("获取服务地址失败: %v", err)
	}

	// 启动监听协程（监听单个键）
	go c.watcher(serviceName)
	return addr, nil
}

// 监听单个键的变更（不再使用前缀）
func (c *EtcdClient) watcher(serviceName string) {
	watchChan := c.watch.Watch(context.Background(), serviceName)
	for resp := range watchChan {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.PUT:
				addr := string(event.Kv.Value)
				c.SetServiceList(serviceName, addr)
				fmt.Printf("服务 %s 地址更新为: %s\n", serviceName, addr)
			case mvccpb.DELETE:
				c.DelServiceList(serviceName)
				fmt.Printf("服务 %s 已下线\n", serviceName)
			}
		}
	}
}
