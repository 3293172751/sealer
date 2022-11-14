package kubernets

import (
	"fmt"
	"sealer-runtime-demo/runtime"
	"sync"
)

type KubeadmRuntime struct {
	*sync.Mutex        // 表示这个结构体是一个锁
	Cluster     string `json:"cluster"` //cluster表示集群的名字
	Action      string `json:"action"`  //action表示集群的动作
	Config      string `json:"config"`  //config表示集群的配置文件
}

func (k *KubeadmRuntime) Reset() error {
	fmt.Println("will reset the k8s...")
	return nil
}

func (k *KubeadmRuntime) Init() error {
	fmt.Println("k8s init start...")
	return k.init()
}

func (k *KubeadmRuntime) Upgrade() error {
	fmt.Println("k8s upgrade start...")
	return nil
}

func (k *KubeadmRuntime) GetMetadata() (string, error) {
	fmt.Println("k8s get metadata...")
	return "", nil
}

func (k *KubeadmRuntime) UpdateCert(cert []string) error {
	fmt.Println("k8s update cert...")
	return nil
}

func NewK8sRuntime() runtime.Interface {
	fmt.Println("judge k8s new runtime...")
	k8s := &KubeadmRuntime{}
	// 打印结构体json
	fmt.Printf("k8s json: %v \r	", k8s)

	return k8s
}
