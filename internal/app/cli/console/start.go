package console

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/derain/internal/pkg/vars"
	"strings"
)

const (
	SYS_COMMAND  = "sys"
	FSYS_COMMAND = "fsys"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "sys", Description: "system information"},
		{Text: "fsys", Description: "file system information"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Start() {
	for {
		t := prompt.Input("> ", completer)
		if len(strings.TrimSpace(t)) == 0 {
			continue
		}
		execCommand(t)
	}
}

func execCommand(t string) {
	switch t {
	case SYS_COMMAND:
		{
			s, _ := json.MarshalIndent(vars.TSys, "", " ")
			fmt.Println(string(s))
			break
		}
	case FSYS_COMMAND:
		{
			s, _ := json.MarshalIndent(vars.TFSys, "", " ")
			fmt.Println(string(s))
			break
		}
	default:
		{
			fmt.Println("non-existent script")
		}
	}
}
