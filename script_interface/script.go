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
	"context"
	"github.com/Shopify/go-lua"
	"github.com/markbates/pkger"
	"github.com/pujo-j/luabox"
	"github.com/pujo-j/luabox/localenv"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
)

func NewEnv(args []string, log *zap.SugaredLogger, libFolder string) (*luabox.Environment, error) {
	fs := luabox.VFS{}
	fs.BaseFs = &localenv.Fs{BaseDir: path.Clean(libFolder)}
	fs.Prefixes = map[string]luabox.Filesystem{}
	env := make(map[string]string)
	for _, k := range os.Environ() {
		env[k] = os.Getenv(k)
	}
	libFile, err := pkger.Open("/lua/ocilot_init.lua")
	if err != nil {
		return nil, err
	}
	defer func() { _ = libFile.Close() }()
	libSource, err := ioutil.ReadAll(libFile)
	if err != nil {
		return nil, err
	}
	lib := luabox.LuaFile{
		Name: "ocilot",
		Code: string(libSource),
	}
	libs := make(map[string]luabox.LuaFile)
	for s, file := range luabox.BaseLibs {
		libs[s] = file
	}
	if libFolder != "" {

	}
	golibs := []lua.RegistryFunction{
		{"ocisys", OciSysOpen},
	}
	res := luabox.Environment{
		Fs:      &fs,
		Context: context.Background(),
		Args:    args,
		Env:     env,
		Input:   os.Stdin,
		Output:  os.Stdout,
		Log: &localenv.ZapLog{
			Zap: log,
		},
		GoLibs:     golibs,
		PreInitLua: []luabox.LuaFile{lib},
		LuaLibs:    libs,
	}
	return &res, nil
}
