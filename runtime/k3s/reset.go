package k3s

import (
	"fmt"
)

func (k3s *K3sRuntime) reset() error {
	fmt.Println("k3s reset ...")
	return nil
}
