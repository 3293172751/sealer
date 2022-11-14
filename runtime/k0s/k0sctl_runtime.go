package k0s

import (
	"fmt"
	"sealer-runtime-demo/runtime"
)

type K0sctlRuntime struct {
	Cluster string `json:"cluster"`
	Action  string `json:"action"`
	Config  string `json:"config"`
}

func (k *K0sctlRuntime) Init() error {
	fmt.Println("k0s init start...")
	return nil
}

func (k *K0sctlRuntime) Reset() error {
	fmt.Println("will reset the k0s...")
	return nil
}

func (k *K0sctlRuntime) Upgrade() error {
	fmt.Println("k0s upgrade start...")
	return nil
}

func (k *K0sctlRuntime) GetMetadata() (string, error) {
	fmt.Println("k0s get metadata...")
	return "", nil
}

func (k *K0sctlRuntime) UpdateCert(cert []string) error {
	fmt.Println("k0s update cert...")
	return nil
}

func NewK0sctlRuntime() runtime.Interface {
	fmt.Println("judge k0s new runtime...")
	return &K0sctlRuntime{}
}
