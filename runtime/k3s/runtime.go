package

import (
	"fmt"
	"sealer-runtime-demo/runtime"
	
)



type k3sRuntime struct {
	ClusterRuntime string `json:"clusterRuntime"`
	Action		 string `json:"action"`
	Config		 string `json:"config"`
}

func (k *k3sRuntime) Init() error {
	fmt.Println("k3s init start...")
	return nil
}

func (k *k3sRuntime) Upgrade() error {
	fmt.Println("k3s upgrade start...")
	return nil
}

func (k *k3sRuntime) Reset() error {
	fmt.Println("k3s reset start...")
	return nil
}

