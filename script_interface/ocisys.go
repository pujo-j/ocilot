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
	_ "github.com/pujo-j/luabox"
)

var ocisys = []lua.RegistryFunction{

	{"snap", LuaSnapShot},
	{"snapDiff", LuaDiff},
	{"snapAsLayer", LuaSnapAsLayer},
	{"snapFiles", LuaSnapFiles},
	{"snapString", LuaSnapString},

	{"layerInfo", LuaLayerInfo},
	{"layerString", LuaLayerString},

	{"imagePull", LuaPullImage},
	{"imagePush", LuaImagePush},
	{"imageClone", LuaImageClone},
	{"imageGetConfig", LuaImageGetConfig},
	{"imageSetConfig", LuaImageSetConfig},
	{"imageGetLayers", LuaImageGetLayers},
	{"imageString", LuaImageString},
	{"imageAppendLayer", LuaImageAppendLayer},

	{"cacheGetImage", LuaGetImageFromCache},

	{"shellExec", LuaShellExec},
	{"shellString", LuaShellString},

	{"hash", LuaHash},
	{"read", LuaRead},
	{"write", LuaWrite},
}

func OciSysOpen(l *lua.State) int {
	lua.NewLibrary(l, ocisys)
	return 1
}
