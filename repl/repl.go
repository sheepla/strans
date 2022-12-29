package repl

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/peterh/liner"
	"github.com/sheepla/strans/api"
)

func promptString(param *api.Param) string {
	return fmt.Sprintf("[%s] %s -> %s\n > ",
		param.Engine,
		param.SourceLang,
		param.TargetLang,
	)
}

func Start(param *api.Param) {
	//nolint:forbidigo
	fmt.Println("REPL mode. Type Ctrl-D to exit.")

	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
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

		param.Text = input
		result, err := api.Translate(param)

		fmt.Fprintln(os.Stdout, result.Text)
	}
}
