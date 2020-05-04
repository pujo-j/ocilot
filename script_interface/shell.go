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
	"fmt"
	"github.com/Shopify/go-lua"
	"os"
	"os/exec"
)

func internalShell(l *lua.State) (*exec.Cmd, error) {
	shell := lua.CheckString(l, 1)
	var args = make([]string, 0, 2)
	args = append(args, lua.CheckString(l, 2))
	cmd := lua.CheckString(l, 3)
	args = append(args, cmd)
	command := exec.Command(shell, args...)
	return command, nil
}

func LuaShellExec(l *lua.State) int {
	cmd, err := internalShell(l)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("executing: %v\n", cmd)
	err = cmd.Run()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	return 0
}

func LuaShellString(l *lua.State) int {
	cmd, err := internalShell(l)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	output, err := cmd.Output()
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushString(string(output))
	return 1
}
