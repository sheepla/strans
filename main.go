//nolint:gocritic,deadcode,varnamelen
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sheepla/strans/api"
	"github.com/sheepla/strans/repl"
	"github.com/urfave/cli/v2"
)

//nolint:gochecknoglobals
var (
	appName        = "strans"
	appDescription = "a command line translate tool with GNU Readline like interactive mode"
	appVersion     = "unknown"
	appRevision    = "unknown"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrAPI
	exitCodeErrIO
	exitCodeErrInternal
)

func (e exitCode) Int() int {
	return int(e)
}

//nolint:gochecknoglobals
var selectedEngine api.Engine

func main() {
	if err := initApp().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[ ERR ] %s\n", err)
	}
}

//nolint:exhaustivestruct,exhaustruct,funlen
func initApp() *cli.App {
	app := &cli.App{
		Name:      appName,
		Usage:     appDescription,
		ArgsUsage: "TEXT...",
		Version:   fmt.Sprintf("%s-%s", appVersion, appRevision),
		UsageText: fmt.Sprintf("%s [OPTIONS] TEXT...\necho TEXT... | %s [OPTIONS] -\n%s [OPTIONS] - < FILE", appName, appName, appName),
		Suggest:   true,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "source",
			Aliases:  []string{"s", "from"},
			Required: false,
			Usage:    "Source language to translate",
			EnvVars:  []string{"STRANS_SOURCE_LANG"},
			Action: func(ctx *cli.Context, s string) error {
				if strings.TrimSpace(s) == "" {
					return cli.Exit(
						"source language must not be empty string",
						exitCodeErrArgs.Int(),
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
						exitCodeErrArgs.Int(),
					)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:     "engine",
			Aliases:  []string{"e"},
			Required: false,
			Usage:    "Name of translate engine",
			EnvVars:  []string{"STRANS_ENGINE"},
			Action: func(ctx *cli.Context, s string) error {
				if strings.TrimSpace(s) == "" {
					return cli.Exit(
						"engine must not be empty string",
						exitCodeErrArgs.Int(),
					)
				}

				eng, err := api.ParseEngineString(s)
				if err != nil {
					return cli.Exit(
						err,
						exitCodeErrArgs.Int(),
					)
				}

				selectedEngine = eng

				return nil
			},
		},
		&cli.StringFlag{
			Name:     "instance",
			Aliases:  []string{"i"},
			Required: false,
			Usage:    "Instance host name of SimplyTranslate",
			EnvVars:  []string{"STRANS_INSTANCE"},
			Action: func(ctx *cli.Context, s string) error {
				if strings.TrimSpace(s) == "" {
					return cli.Exit(
						"instance must not be empty string",
						exitCodeErrArgs.Int(),
					)
				}

				return nil
			},
		},
		&cli.BoolFlag{
			Name:     "repl",
			Aliases:  []string{"r"},
			Required: false,
			Usage:    "Start interactive mode",
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
	source := ctx.String("source")
	target := ctx.String("target")
	instance := ctx.String("instance")

	var text string

	if ctx.NArg() == 1 && ctx.Args().First() == "-" {
		var err error

		text, err = readString(ctx.App.Reader)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf("failed to read stnadard input: %s", err),
				exitCodeErrIO.Int(),
			)
		}
	} else {
		text = strings.Join(ctx.Args().Slice(), " ")
	}

	// Create parameter
	param, err := api.NewParam(source, target, text, selectedEngine, instance)
	if err != nil {
		return cli.Exit(
			err,
			exitCodeErrArgs.Int(),
		)
	}

	if ctx.Bool("repl") {
		// Start REPL mode
		repl.Start(param)

		return cli.Exit("", exitCodeOK.Int())
	}

	// Execute translate
	result, err := api.Translate(param)
	if err != nil {
		if errors.Is(err, api.ErrAPI) {
			return cli.Exit(
				fmt.Sprintf("%s: %s", api.ErrAPI, err),
				exitCodeErrAPI.Int(),
			)
		}

		return cli.Exit(
			fmt.Sprintf("internal error: %s", err),
			exitCodeErrInternal.Int(),
		)
	}

	fmt.Fprintln(ctx.App.Writer, result.Text)

	return cli.Exit("", exitCodeOK.Int())
}

func readString(r io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		//nolint:wrapcheck
		return "", err
	}

	return buf.String(), nil
}
