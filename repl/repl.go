package repl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
	"github.com/sheepla/strans/audio"
	"github.com/sheepla/strans/trans"
)

// func promptString(param *api.Param) string {
// 	return fmt.Sprintf("[%s] %s -> %s\n > ",
// 		param.Engine,
// 		param.SourceLang,
// 		param.TargetLang,
// 	)
// }

func Start(param *trans.TranslateParam, playAudio bool) {
	//nolint:forbidigo
	fmt.Println("REPL mode. Type Ctrl-D to exit.")

	line := liner.NewLiner()
	line.SetCtrlCAborts(false)
	line.SetMultiLineMode(true)

	defer line.Close()

	for {
		input, err := line.Prompt("> ")
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		param.Text = input

		result, err := trans.Translate(param)
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
