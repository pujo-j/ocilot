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
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"strings"
)

func GetImageFromCache(baseUrl string, key string) (*Image, bool, error) {
	if strings.HasPrefix(baseUrl, "docker://") {
		return nil, false, errors.New("cache needs a remote registry, not a local docker storage")
	}
	h := sha256.New()
	h.Write([]byte(key))
	ref := baseUrl + hex.EncodeToString(h.Sum(nil))
	image, err := LoadImage(ref)
	if err != nil {
		tmp := Image{
			ref:      nil,
			img:      empty.Image,
			isDocker: false,
		}
		clone, err := tmp.Clone(ref)
		if err != nil {
			return nil, false, err
		}
		return clone, false, nil
	}
	return image, true, nil
}
