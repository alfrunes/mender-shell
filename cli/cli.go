// Copyright 2020 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package cli

import (
	"github.com/urfave/cli/v2"

	"github.com/mendersoftware/mender-shell/config"
)

func SetupCLI(args []string) error {
	runOptions := &runOptionsType{}
	app := &cli.App{
		Description: "",
		Name:        "mender-shell",
		Usage:       "manage and start the Mender shell.",
		Version:     config.ShowVersion(),
		Commands: []*cli.Command{
			{
				Name:   "daemon",
				Usage:  "Start the client as a background service.",
				Action: runOptions.handleCLIOptions,
			},
			{
				Name:   "version",
				Usage:  "Show the version and runtime information of the binary build",
				Action: config.ShowVersionCLI,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Configuration `FILE` path.",
				Value:       config.DefaultConfFile,
				Destination: &runOptions.config,
			},
			&cli.StringFlag{
				Name:        "fallback-config",
				Aliases:     []string{"b"},
				Usage:       "Fallback configuration `FILE` path.",
				Value:       config.DefaultFallbackConfFile,
				Destination: &runOptions.fallbackConfig,
			},
		},
	}

	return app.Run(args)
}

func (runOptions *runOptionsType) handleCLIOptions(ctx *cli.Context) error {
	// Handle config flags
	config, err := config.LoadConfig(runOptions.config, runOptions.fallbackConfig)
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return err
	}

	switch ctx.Command.Name {
	case "daemon":
		d, err := initDaemon(config)
		if err != nil {
			return err
		}
		return runDaemon(d)
	default:
		cli.ShowAppHelpAndExit(ctx, 1)
	}
	return nil
}
