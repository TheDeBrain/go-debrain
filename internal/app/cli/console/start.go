package console

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/derain/core/db/table/sys"
	"os"
	"strings"
)

const (
	SYS_COMMAND  = "sys"
	FSYS_COMMAND = "fsys"
	EXIT_COMMAND = "exit"
	QUIT_COMMAND = "quit"
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
			s, _ := json.MarshalIndent(sys.LoadTSys(), "", " ")
			fmt.Println(string(s))
			break
		}
	case FSYS_COMMAND:
		{
			s, _ := json.MarshalIndent(sys.LoadFileSys(), "", " ")
			fmt.Println(string(s))
			break
		}
	case EXIT_COMMAND:
		{
			os.Exit(0)
			break
		}
	case QUIT_COMMAND:
		{
			os.Exit(0)
			break
		}
	default:
		{
			fmt.Println("non-existent script")
		}
	}
}
