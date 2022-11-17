package k3s

import (
	"fmt"

	"github.com/cubxxw/sealer-runtime/runtime"
)

type K3sRuntime struct {
	ClusterRuntime string `json:"clusterRuntime"`
	Action         string `json:"action"`
	Config         string `json:"config"`
}

func (k *K3sRuntime) Init() error {
	fmt.Println("K3sRuntime start to create a cluster ...")
	return k.init()
}

func (k *K3sRuntime) Upgrade() error {
	fmt.Println("K3sRuntime start to upgrade a cluster ...")
	return nil
}

func (k *K3sRuntime) Reset() error {
	fmt.Println("K3sRuntime start to reset a cluster ...")
	return nil
}

func (k *K3sRuntime) GetMetadata() (string, error) {
	fmt.Println("K3sRuntime start to get metadata ...")
	return "K3sRuntime", nil
}

func (k *K3sRuntime) UpdateCert(certs []string) error {
	fmt.Println("K3sRuntime start to update certs ...")
	return nil
}

func NewK3sRuntime() runtime.Interface {
	fmt.Println("judge k8s new runtime...")
	k3s := &K3sRuntime{}
	// 打印结构体json
	fmt.Printf("k8s json: %v \r	", k3s)

	return k3s
}
