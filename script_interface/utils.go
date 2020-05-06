package script_interface

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Shopify/go-lua"
	"io/ioutil"
	"os"
)

func LuaHash(l *lua.State) int {
	data := lua.CheckString(l, 1)
	h := sha256.New()
	h.Write([]byte(data))
	res := hex.EncodeToString(h.Sum(nil))
	l.PushString(res)
	return 1
}

func LuaRead(l *lua.State) int {
	fileName := lua.CheckString(l, 1)
	stat, err := os.Stat(fileName)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	if stat.IsDir() {
		l.PushString(fileName + " is a directory")
		l.Error()
		return 0
	}
	if stat.Size() > 256*1024 {
		l.PushString("refusing to load >256Kib in memory")
		l.Error()
		return 0
	}
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushString(string(bytes))
	return 1
}

func LuaWrite(l *lua.State) int {
	fileName := lua.CheckString(l, 1)
	data := lua.CheckString(l, 2)
	err := ioutil.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	return 0
}
