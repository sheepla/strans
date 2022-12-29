//nolint:gocritic,deadcode
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sheepla/strans/api"
	"github.com/urfave/cli/v2"
)

//nolint:gochecknoglobals
var (
	appName        = "strans"
	appDescription = "a command line SimplyTranslate client with bash-like interactive mode"
	appVersion     = "unknown"
)

type exitCode int

const (
	exitCodeOK = iota
	exitCodeErrArgs
	exitCodeErrAPI
	exitCodeErrInternal
)

func main() {
	if err := initApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

//nolint:exhaustivestruct,exhaustruct,funlen
func initApp() *cli.App {
	app := &cli.App{
		Name:      appName,
		Usage:     appDescription,
		ArgsUsage: "TEXT...",
		Version:   appVersion,
		UsageText: fmt.Sprintf("%s [-e|--engine ENGINE][-i|--instance INSTANCE] -s SOURCE_LANG -t TARGET_LANG TEXT...\n%s [-r|--repl]", appName, appName),
		Suggest:   true,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "source",
			Aliases:  []string{"s", "from"},
			Required: true,
			Usage:    "Source language to translate",
			EnvVars:  []string{"STRANS_SOURCE_LANG"},
			Action: func(ctx *cli.Context, s string) error {
				if strings.TrimSpace(s) == "" {
					return cli.Exit(
						"source language must not be empty string",
						exitCodeErrArgs,
					)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:     "target",
			Aliases:  []string{"t", "to"},
			Required: true,
			Usage:    "Target language to translate",
			EnvVars:  []string{"STRANS_TARGET_LANG"},
			Action: func(ctx *cli.Context, s string) error {
				if strings.TrimSpace(s) == "" {
					return cli.Exit(
						"target language must not be empty string",
						exitCodeErrArgs,
					)
				}

				return nil
			},
		},
		//&cli.StringFlag{
		//	Name:     "engine",
		//	Aliases:  []string{"e"},
		//	Required: false,
		//	Usage:    "Name of translate engine",
		//	Action: func(ctx *cli.Context, s string) error {
		//		if strings.TrimSpace(s) == "" {
		//			return cli.Exit(
		//				"target language must not be empty string",
		//				exitCodeErrArgs,
		//			)
		//		}

		//		return nil
		//	},
		//},
		&cli.StringFlag{
			Name:     "instance",
			Aliases:  []string{"i"},
			Required: false,
			Usage:    "Instance URL of SimplyTranslate",
			EnvVars:  []string{"STRANS_INSTANCE"},
			Action: func(ctx *cli.Context, s string) error {
				if strings.TrimSpace(s) == "" {
					return cli.Exit(
						"target language must not be empty string",
						exitCodeErrArgs,
					)
				}

				return nil
			},
		},
		&cli.BoolFlag{
			Name:     "repl",
			Aliases:  []string{"r"},
			Required: false,
			Usage:    "Start bash-like REPL mode",
		},
		&cli.BoolFlag{
			Name:     "debug",
			Aliases:  []string{},
			Required: false,
			Usage:    "Enable debug mode",
			Hidden:   true,
		},
	}

	app.Action = run

	return app
}

func run(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return cli.Exit(
			"must require argument(s)",
			exitCodeErrArgs,
		)
	}

	source := ctx.String("source")
	target := ctx.String("target")

	//engine, err := api.ParseEngineString(ctx.String("engine"))
	//if err != nil {
	//	return cli.Exit(
	//		err,
	//		exitCodeErrArgs,
	//	)
	//}

	text := strings.Join(ctx.Args().Slice(), " ")

	param, err := api.NewParam(source, target, text, api.EngineDefault)
	if err != nil {
		return cli.Exit(
			err,
			exitCodeErrArgs,
		)
	}

	instance := ctx.String("instance")

	result, err := api.Translate(param, instance)
	if err != nil {
		if errors.Is(err, api.ErrAPI) {
			return cli.Exit(
				fmt.Sprintf("%s: %s", api.ErrAPI, err),
				exitCodeErrAPI,
			)
		}

		return cli.Exit(
			fmt.Sprintf("internal error: %s", err),
			exitCodeErrInternal,
		)
	}

	fmt.Fprintln(ctx.App.Writer, result.Text)

	return cli.Exit("", exitCodeOK)
}
