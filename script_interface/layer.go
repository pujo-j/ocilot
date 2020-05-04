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

package script_interface

import (
	"github.com/Shopify/go-lua"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pujo-j/luabox"
)

func LuaLayerString(l *lua.State) int {
	lua.CheckAny(l, 1)
	li := l.ToUserData(1)
	layer, ok := li.(v1.Layer)
	if !ok {
		l.PushString("Expected layer as parameter")
		l.Error()
		return 0
	}
	digest, err := layer.Digest()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushString(digest.Hex)
	return 1
}

func LuaLayerInfo(l *lua.State) int {
	lua.CheckAny(l, 1)
	li := l.ToUserData(1)
	layer, ok := li.(v1.Layer)
	if !ok {
		l.PushString("Expected layer as parameter")
		l.Error()
		return 0
	}
	var err error
	res := make(map[string]interface{})
	diffid, err := layer.DiffID()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	res["diffId"] = diffid.Hex
	digest, err := layer.Digest()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	res["digest"] = digest.Hex
	size, err := layer.Size()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	res["size"] = size
	mediaType, err := layer.MediaType()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	res["mime"] = string(mediaType)
	luabox.DeepPush(l, res)
	return 1
}
