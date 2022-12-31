package repl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	"github.com/sheepla/strans/api"
	"github.com/sheepla/strans/audio"
)

//nolint:gochecknoglobals
var historyFileName = filepath.Join(os.TempDir(), "strans_history.txt")

//nolint:cyclop,funlen
func Start(param *api.TranslateParam, playAudio bool) {
	fmt.Fprintln(os.Stdout, "Interactive mode. Type Ctrl-D to exit.")

	line := liner.NewLiner()
	line.SetCtrlCAborts(false)
	line.SetMultiLineMode(true)

	defer line.Close()

	//nolint:gomnd,nosnakecase
	historyFile, err := os.OpenFile(historyFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	_, err = line.ReadHistory(historyFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	defer historyFile.Close()

	defer func() {
		if _, err := line.WriteHistory(historyFile); err != nil {
			fmt.Fprintln(os.Stdout, err)
		}
	}()

REPL:
	for {
		fmt.Fprintln(os.Stdout, newPrompt(param))

		input, err := line.Prompt("> ")
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(os.Stdout, "bye")

				break REPL
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		param.Text = input

		line.AppendHistory(input)

		result, err := api.Translate(param)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Fprintln(os.Stdout, result.Text)

		if playAudio {
			if err := audio.FetchAndPlay(param.TargetLang, result.Text, param.Instance); err != nil {
				fmt.Fprintln(os.Stdout, err)
			}
		}
	}
}

func newPrompt(param *api.TranslateParam) string {
	return fmt.Sprintf("\n[%s -> %s]", param.SourceLang, param.TargetLang)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
