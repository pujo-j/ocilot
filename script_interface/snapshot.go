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
	"github.com/pujo-j/luabox"
	"ocilot"
)

func LuaSnapString(l *lua.State) int {
	lua.CheckAny(l, 1)
	s1i := l.ToUserData(1)
	s1, ok := s1i.(*ocilot.Snapshot)
	if !ok {
		l.PushString("Expected snapshot as parameter")
		l.Error()
		return 0
	}
	l.PushString(s1.String())
	return 1
}

func LuaSnapShot(l *lua.State) int {
	filePath := lua.CheckString(l, 1)
	snapshot, err := ocilot.NewSnapshot(filePath)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushUserData(snapshot)
	return 1
}

func LuaDiff(l *lua.State) int {
	lua.CheckAny(l, 1)
	lua.CheckAny(l, 2)
	s1i := l.ToUserData(1)
	s1, ok := s1i.(*ocilot.Snapshot)
	if !ok {
		l.PushString("Expected snapshot as first parameter")
		l.Error()
		return 0
	}
	s2i := l.ToUserData(2)
	s2, ok := s2i.(*ocilot.Snapshot)
	if !ok {
		l.PushString("Expected snapshot as first parameter")
		l.Error()
		return 0
	}
	diff, err := ocilot.Diff(s1, s2)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushUserData(diff)
	return 1
}

func LuaSnapAsLayer(l *lua.State) int {
	lua.CheckAny(l, 1)
	workFile := lua.CheckString(l, 2)
	s1i := l.ToUserData(1)
	s1, ok := s1i.(*ocilot.Snapshot)
	if !ok {
		l.PushString("Expected snapshot as parameter")
		l.Error()
		return 0
	}
	layer, err := s1.AsLayer(workFile)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushUserData(layer)
	return 1
}

func LuaSnapFiles(l *lua.State) int {
	lua.CheckAny(l, 1)
	s1i := l.ToUserData(1)
	s1, ok := s1i.(*ocilot.Snapshot)
	if !ok {
		l.PushString("Expected snapshot as parameter")
		l.Error()
		return 0
	}
	res := make([]string, 0, len(s1.Files))
	for _, header := range s1.Files {
		res = append(res, header.Name)
	}
	luabox.DeepPush(l, res)
	return 1
}
