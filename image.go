/*
 *    Copyright 2020 Josselin Pujo
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package ocilot

import (
	"encoding/json"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"strings"
)

type Image struct {
	ref      name.Reference
	img      v1.Image
	isDocker bool
}

var keyChain = authn.NewMultiKeychain(authn.DefaultKeychain, google.Keychain)

func LoadImage(imageName string) (*Image, error) {
	if strings.HasPrefix(imageName, "docker://") {
		reference, err := name.ParseReference(strings.TrimPrefix(imageName, "docker://"))
		if err != nil {
			return nil, err
		}
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
		image, err := daemon.Image(reference, daemon.WithClient(dockerClient))
		if err != nil {
			return nil, err
		}
		return &Image{
			ref:      reference,
			img:      image,
			isDocker: true,
		}, nil
	} else {
		reference, err := name.ParseReference(imageName)
		if err != nil {
			return nil, err
		}
		descriptor, err := remote.Get(reference, remote.WithAuthFromKeychain(keyChain))
		if err != nil {
			return nil, err
		}
		img, err := descriptor.Image()
		return &Image{
			ref:      reference,
			img:      img,
			isDocker: false,
		}, nil
	}
}

func (i *Image) String() string {
	return i.ref.String()
}

func (i *Image) Clone(targetName string) (*Image, error) {
	if strings.HasPrefix(targetName, "docker://") {
		reference, err := name.ParseReference(strings.TrimPrefix(targetName, "docker://"))
		if err != nil {
			return nil, err
		}
		return &Image{
			ref:      reference,
			img:      i.img,
			isDocker: true,
		}, nil
	} else {
		reference, err := name.ParseReference(targetName)
		if err != nil {
			return nil, err
		}
		return &Image{
			ref:      reference,
			img:      i.img,
			isDocker: false,
		}, nil
	}

}

func (i *Image) Push() error {
	if i.isDocker {
		tag, err := name.NewTag(i.ref.Name())
		if err != nil {
			return err
		}
		_, err = daemon.Write(tag, i.img)
		return err
	} else {
		return remote.Write(i.ref, i.img, remote.WithAuthFromKeychain(keyChain))
	}
}

func (i *Image) AddLayer(layer v1.Layer) error {
	image, err := mutate.AppendLayers(i.img, layer)
	if err != nil {
		return err
	}
	i.img = image
	return nil
}

func (i *Image) GetConfig() (*v1.Config, error) {
	configFile, err := i.img.ConfigFile()
	if err != nil {
		return nil, err
	}
	// Deep copy
	marshal, err := json.Marshal(configFile.Config)
	res := v1.Config{}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *Image) WithConfig(config *v1.Config) error {
	configFile, err := i.img.ConfigFile()
	if err != nil {
		return err
	}
	configFile.Config = *config
	image, err := mutate.ConfigFile(i.img, configFile)
	if err != nil {
		return err
	}
	i.img = image
	return nil
}

func (i *Image) Layers() ([]v1.Layer, error) {
	return i.img.Layers()
}
