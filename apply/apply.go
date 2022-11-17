/*
Here's what the selected code is doing:
1.Import the runtime package
2. Create a Processor struct
3. Create a method on the Processor struct called ApplyClusterFile
4. Create a switch statement that will set the runtime interface to the correct runtime
5. Create a method on the Processor struct called GetMetadata
6. Create a method on the Processor struct called Init

Now that we have the code, let's run it:

	go test .\runtime\kubernets\
*/
package apply

import (
	"fmt"

	"github.com/cubxxw/sealer-runtime/runtime"
	"github.com/cubxxw/sealer-runtime/runtime/k0s"
	"github.com/cubxxw/sealer-runtime/runtime/k3s"
	"github.com/cubxxw/sealer-runtime/runtime/kubernets"
)

type Processor struct {
	Runtime runtime.Interface
	//runtime.Interface is an interface that is implemented by the runtime package
}

func (c *Processor) ApplyClusterFile() error {
	c.Runtime = runtime.NewKsRuntime(KubernetsRuntime)

	metadata, err := c.Runtime.GetMetadata()
	if err != nil {
		return fmt.Errorf("failed to get runtime metadata: %w", err)
	}

	switch metadata {
	case "k8s":
		c.Runtime = kubernets.NewK8sRuntime()
	case "k0s":
		c.Runtime = k0s.NewK0sctlRuntime()
	case "k3s":
		c.Runtime = k3s.NewK3sRuntime()
	default:
		c.Runtime = kubernets.NewK8sRuntime()
	}

	return c.Runtime.Init()
}

func NewProcessor() *Processor {
	return &Processor{}
}
