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
package main

import (
	"github.com/Shopify/go-lua"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	script "ocilot/script_interface"
	"os"
)

func main() {
	_ = rootCmd.Execute()
}

var log *zap.SugaredLogger
var rootCmd = &cobra.Command{
	Use:     "ocilot",
	Short:   "ocilot provides a scripting language for simple OCI image manipulation",
	Example: "ocilot script.lua",
	Args:    cobra.MinimumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		json, err := cmd.Flags().GetBool("log-json")
		if err != nil {
			panic(err)
		}

		var logconfig zap.Config
		if json {
			logconfig = zap.NewProductionConfig()
		} else {
			logconfig = zap.NewDevelopmentConfig()
		}
		logconfig.DisableCaller = true
		debug, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			panic(err)
		}
		if debug {
			logconfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
			logger, err := logconfig.Build()
			if err != nil {
				panic(err)
			}
			log = logger.Sugar()
		} else {
			logconfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
			logger, err := logconfig.Build()
			if err != nil {
				panic(err)
			}
			log = logger.Sugar()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var remaining []string
		if len(args) > 1 {
			remaining = args[1:]
		} else {
			remaining = []string{}
		}
		libFolder, err := cmd.PersistentFlags().GetString("lib")
		if err != nil {
			return err
		}
		env, err := script.NewEnv(remaining, log, libFolder)
		if err != nil {
			return err
		}
		l, err := env.Init()
		if err != nil {
			return err
		}
		scriptName := args[0]
		if scriptName == "-" {
			err = l.Load(os.Stdin, "-", "text")
			if err != nil {
				log.With("error", err).Error("loading lua script")
				return err
			}
		} else {
			f, err := os.Open(scriptName)
			if err != nil {
				log.With("file", scriptName, "error", err).Error("opening file")
				return err
			}
			defer func() {
				_ = f.Close()
			}()
			err = l.Load(f, "-", "text")
			if err != nil {
				log.With("error", err).Error("loading lua script")
				return err
			}
		}

		err = l.ProtectedCall(0, lua.MultipleReturns, 0)
		if err != nil {
			log.With("error", err).Error("executing lua script")
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose Logging")
	rootCmd.PersistentFlags().BoolP("log-json", "j", false, "JSON logging output")
	rootCmd.PersistentFlags().StringP("lib", "l", "/ocilot", "lua libraries folder")
}
