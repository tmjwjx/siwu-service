package inits

import (
	"forum/pkg/etcd"
	"forum/pkg/globals"
	"log"
)

func InitEtcd(addr []string) {
	etcd, err := etcd.NewEtcdClient(addr)
	if err != nil {
		log.Println("InitEtcd() -> 创建 etcdService 失败, err = ", err)
	}
	globals.EtcdClient = etcd
}
