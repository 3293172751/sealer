package runtime

/*
对kubernetes、k3s、k0s进行统一管理
区分不同的runtime，返回不同的runtime
*/

func NewKsRuntime(ksname string) interface{} {
	// 判断传入的是k3s ,k0s 还是kubernets
	switch ksname {
	case "k3s":
		return &K3s{}
	case "k0s":
		return &K0s{}
	case "kubernetes":
		return &KubeRuntime{}
	default:
		return nil
	}
}
