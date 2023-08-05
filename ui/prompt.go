package ui

import (
	"fmt"

	fzf "github.com/koki-develop/go-fzf"
	"github.com/sheepla/strans/api"
)

func SelectLang(langs []api.Lang) (int, error) {
    f, err := fzf.New(fzf.WithLimit(1))
    if err != nil {
      return 0, err
    }

    indexes, err := f.Find(langs, func(i int) string {
      return fmt.Sprintf("%s [%s]", langs[i].Name, langs[i].Code)
    })

    return indexes[0], err
}

