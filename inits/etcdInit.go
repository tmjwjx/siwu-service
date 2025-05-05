package inits

import (
	"forum/pkg/etcd"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

func InitEtcd() {

	if err := viper.UnmarshalKey("etcd", &globals.AppConfig.Etcd); err != nil {
		globals.Log.Fatalf("Run -> 无法解码为结构: %s", err)
	}

	etcd, err := etcd.NewEtcdClient([]string{globals.AppConfig.Etcd.Endpoints})
	if err != nil {
		log.Println("InitEtcd() -> 创建 etcdService 失败, err = ", err)
	}

	globals.EtcdClient = etcd
}
