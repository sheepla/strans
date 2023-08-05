package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sheepla/strans/api"
	"github.com/sheepla/strans/audio"
	ui "github.com/sheepla/strans/ui"
	cli "github.com/urfave/cli/v2"
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
	exitCodeErrTranslate
	exitCodeErrIO
	exitCodeErrUI
)

func (e exitCode) Int() int {
	return int(e)
}

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
			Required: true,
			Usage:    "Source language to translate; specifying the value '?' will shows a selectable menu",
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
			Usage:    "Target language to translate; specifying the value '?' will shows a selectable menu",
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
			Name:     "instance",
			Aliases:  []string{"i"},
			Required: false,
			Usage:    "Instance host name of Lingva Translate",
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
			Name:     "audio",
			Aliases:  []string{"a"},
			Required: false,
			Usage:    "Read translated text aloud",
		},
		&cli.BoolFlag{
			Name:     "list-source",
			Aliases:  []string{"S"},
			Usage:    "Show a list of source languages",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "list-target",
			Aliases:  []string{"T"},
			Usage:    "Show a list of target languages",
			Required: false,
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
	if ctx.Bool("list-source") {
		langs, err := api.ListSourceLangs(ctx.String("instance"))
		if err != nil {
			return cli.Exit(
				err,
				exitCodeErrTranslate.Int(),
			)
		}

		printLangs(ctx.App.Writer, langs)

		return cli.Exit("", exitCodeOK.Int())
	}

	if ctx.Bool("list-target") {
		langs, err := api.ListTargetLangs(ctx.String("instance"))
		if err != nil {
			return cli.Exit(
				err,
				exitCodeErrTranslate.Int(),
			)
		}

		printLangs(ctx.App.Writer, langs)

		return cli.Exit("", exitCodeOK.Int())
	}

	// Read stdin mode
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

	// Get instance param
	instance := ctx.String("instance")

	// Get source language param
	source := ctx.String("source")
	if source == "?" {
		sourceLangs, err := api.ListSourceLangs(instance)
		if err != nil {
			return cli.Exit(err, exitCodeErrTranslate.Int())
		}

		idx, err := ui.SelectLang(*sourceLangs)
		if err != nil {
			return cli.Exit(err, exitCodeErrUI.Int())
		}

		source = (*sourceLangs)[idx].Code
	}

	// Get target language param
	target := ctx.String("target")
	if target == "?" {
		targetLangs, err := api.ListSourceLangs(instance)
		if err != nil {
			return cli.Exit(err, exitCodeErrTranslate.Int())
		}

		idx, err := ui.SelectLang(*targetLangs)
		if err != nil {
			return cli.Exit(err, exitCodeErrUI.Int())
		}

		target = (*targetLangs)[idx].Code
	}

	// Create parameter
	param, err := api.NewTranslateParam(source, target, text, instance)
	if err != nil {
		return cli.Exit(err, exitCodeErrArgs.Int())
	}

	// Start REPL mode
	if ctx.Bool("repl") {
		ui.Repl(param, ctx.Bool("audio"))

		return cli.Exit("", exitCodeOK.Int())
	}

	// Execute translate
	result, err := api.Translate(param)
	if err != nil {
		return cli.Exit(err, exitCodeErrTranslate.Int())
	}

	// Output translated text
	fmt.Fprintln(ctx.App.Writer, result.Text)

	// Read translated text aloud
	if ctx.Bool("audio") {
		if err := audio.FetchAndPlay(param.TargetLang, result.Text, param.Instance); err != nil {
			return cli.Exit(err, exitCodeErrIO.Int())
		}
	}

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

func printLangs(w io.Writer, langs *[]api.Lang) {
	for _, lang := range *langs {
		fmt.Fprintf(w, "%s\t%s\n", lang.Code, lang.Name)
	}
}
