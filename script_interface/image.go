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
	"encoding/json"
	"github.com/Shopify/go-lua"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/pujo-j/luabox"
	"ocilot"
)

func LuaPullImage(l *lua.State) int {
	name := lua.CheckString(l, 1)
	newImage, err := ocilot.LoadImage(name)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushUserData(newImage)
	return 1
}

func LuaImageClone(l *lua.State) int {
	lua.CheckAny(l, 1)
	i := l.ToUserData(1)
	name := lua.CheckString(l, 2)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as parameter")
		l.Error()
		return 0
	}
	clone, err := image.Clone(name)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushUserData(clone)
	return 1
}

func LuaImagePush(l *lua.State) int {
	lua.CheckAny(l, 1)
	i := l.ToUserData(1)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as parameter")
		l.Error()
		return 0
	}
	err := image.Push()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	return 0
}

func LuaImageGetConfig(l *lua.State) int {
	lua.CheckAny(l, 1)
	i := l.ToUserData(1)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as first parameter")
		l.Error()
		return 0
	}
	config, err := image.GetConfig()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	res := make(map[string]interface{})
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	luabox.DeepPush(l, res)
	return 1
}

func LuaImageSetConfig(l *lua.State) int {
	lua.CheckAny(l, 1)
	i := l.ToUserData(1)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as first parameter")
		l.Error()
		return 0
	}
	p, err := luabox.PullTable(l, 2)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	data, err := json.Marshal(p)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	config := &v1.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	err = image.WithConfig(config)
	return 0
}

func LuaImageGetLayers(l *lua.State) int {
	lua.CheckAny(l, 1)
	i := l.ToUserData(1)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as first parameter")
		l.Error()
		return 0
	}
	layers, err := image.Layers()
	if err != nil {

		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.CreateTable(len(layers), 0)
	for i2, layer := range layers {
		l.PushUserData(layer)
		l.RawSetInt(-2, i2+1)
	}
	return 1
}

func LuaImageString(l *lua.State) int {
	lua.CheckAny(l, 1)
	i := l.ToUserData(1)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as first parameter")
		l.Error()
		return 0
	}
	l.PushString(image.String())
	return 1
}

func LuaImageAppendLayer(l *lua.State) int {
	lua.CheckAny(l, 1)
	lua.CheckAny(l, 2)
	i := l.ToUserData(1)
	image, ok := i.(*ocilot.Image)
	if !ok {
		l.PushString("Expected image as first parameter")
		l.Error()
		return 0
	}
	li := l.ToUserData(2)
	layer, ok := li.(v1.Layer)
	if !ok {
		l.PushString("Expected layer as second parameter")
		l.Error()
		return 0
	}
	err := image.AddLayer(layer)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	return 0
}
