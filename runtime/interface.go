/*
ApplyClusterFile - Apply the cluster file

Parameters:
	None

Returns:
	error - Any error that may have occurred
*/
/*
ApplyClusterFile - Apply the cluster file

Parameters:
	None

Returns:
	error - Any error that may have occurred
*/
package runtime

import (
	"fmt"
)

type Interface interface {
	Init() error
	Upgrade() error
	Reset() error
	GetMetadata() (string, error)
	UpdateCert(certs []string) error //uodate cert
}

// ClusterRuntime is a flag to distinguish the runtime for k0s、k8s、k3s
type ClusterRuntime string

type Metadata struct {
	Version string `json:"version"`
	Arch    string `json:"arch"`
	Variant string `json:"variant"`
	//KubeVersion is a SemVer constraint specifying the version of Kubernetes required.
	KubeVersion string `json:"kubeVersion"`
	NydusFlag   bool   `json:"NydusFlag"`
}

type KubeRuntime struct {
}

func (k *KubeRuntime) Init() error {
	return nil
}

func (k *KubeRuntime) Upgrade() error {
	return nil
}

func (k *KubeRuntime) Reset() error {
	return nil
}

func (k *KubeRuntime) GetMetadata() (string, error) {
	return "k8s", nil
}

func (k *KubeRuntime) UpdateCert(certs []string) error {
	return nil
}

type K3s struct {
}

func (k *K3s) Init() error {
	fmt.Println("k3s start to create a cluster ...")
	return nil
}

func (k *K3s) Upgrade() error {
	fmt.Println("k3s start to upgrade a cluster ...")
	return nil
}

func (k *K3s) Reset() error {
	fmt.Println("k3s start to reset a cluster ...")
	return nil
}

func (k *K3s) GetMetadata() (string, error) {
	fmt.Println("k3s start to get metadata ...")
	return "k3s", nil
}

func (k *K3s) UpdateCert(certs []string) error {
	fmt.Println("k3s start to update certs ...")
	return nil
}

// k0s inherits k3s
type K0s struct {
}

func (k *K0s) Init() error {
	fmt.Println("k0s start to create a cluster ...")
	return nil
}

func (k *K0s) Upgrade() error {
	fmt.Println("k0s start to upgrade a cluster ...")
	return nil
}

func (k *K0s) Reset() error {
	fmt.Println("k0s start to reset a cluster ...")
	return nil
}

func (k *K0s) GetMetadata() (string, error) {
	fmt.Println("k0s start to get metadata ...")
	return "k0s", nil
}

func (k *K0s) UpdateCert(certs []string) error {
	fmt.Println("k0s start to update certs ...")
	return nil
}
