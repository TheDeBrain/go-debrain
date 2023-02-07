package command

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	SyncPort   string
	RpcPort    string
	WebApiPort string
}

func (c *Command) Start() {
	flag.Usage = usage
	flag.StringVar(&c.SyncPort, "syncport", "", "sync port")
	flag.StringVar(&c.RpcPort, "rpcport", "", "rpc port")
	flag.StringVar(&c.WebApiPort, "webapiport", "", "web api port")
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stdout, "debrain version: debrain/0.0.1 Usage: \n"+
		"debrain [-help] [-syncport port] [-rpcport port] [-webapiport port]\n"+
		"Options: \n")
}
