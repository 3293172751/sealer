// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudimage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/pkg/image"
	"github.com/alibaba/sealer/pkg/image/store"
	v2 "github.com/alibaba/sealer/types/api/v2"
	"github.com/alibaba/sealer/utils"
	"github.com/alibaba/sealer/utils/mount"
)

type Interface interface {
	MountImage(cluster *v2.Cluster) error
	UnMountImage(cluster *v2.Cluster) error
}

type mounter struct {
	imageStore store.ImageStore
}

func (p *mounter) MountImage(cluster *v2.Cluster) error {
	return p.mountImage(cluster)
}

func (p *mounter) UnMountImage(cluster *v2.Cluster) error {
	return p.umountImage(cluster)
}

func (p *mounter) umountImage(cluster *v2.Cluster) error {
	mountDir := common.DefaultMountCloudImageDir(cluster.Name)
	if !utils.IsFileExist(mountDir) {
		return nil
	}
	if isMount, _ := mount.GetMountDetails(mountDir); isMount {
		err := utils.Retry(10, time.Second, func() error {
			return mount.NewMountDriver().Unmount(mountDir)
		})
		if err != nil {
			return fmt.Errorf("failed to unmount dir %s,err: %v", mountDir, err)
		}
	}
	return os.RemoveAll(mountDir)
}

func (p *mounter) mountImage(cluster *v2.Cluster) error {
	var (
		mountdir = common.DefaultMountCloudImageDir(cluster.Name)
		upperDir = filepath.Join(mountdir, "upper")
		driver   = mount.NewMountDriver()
		err      error
	)
	if isMount, _ := mount.GetMountDetails(mountdir); isMount {
		err = driver.Unmount(mountdir)
		if err != nil {
			return fmt.Errorf("%s already mount, and failed to umount %v", mountdir, err)
		}
	}
	if utils.IsFileExist(mountdir) {
		err = os.RemoveAll(mountdir)
		if err != nil {
			return fmt.Errorf("failed to clean %s, %v", mountdir, err)
		}
	}
	//get layers
	Image, err := p.imageStore.GetByName(cluster.Spec.Image)
	if err != nil {
		return err
	}
	layers, err := image.GetImageLayerDirs(Image)
	if err != nil {
		return fmt.Errorf("get layers failed: %v", err)
	}

	if err = os.MkdirAll(upperDir, 0744); err != nil {
		return fmt.Errorf("create upperdir failed, %s", err)
	}
	if err = driver.Mount(mountdir, upperDir, layers...); err != nil {
		return fmt.Errorf("mount files failed %v", err)
	}
	return nil
}

func NewCloudImageMounter(is store.ImageStore) (Interface, error) {
	return &mounter{
		imageStore: is,
	}, nil
}