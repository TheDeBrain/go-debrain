package command

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	SyncPort            string `json:"sync_port"`               // sync port,default:9000
	RpcPort             string `json:"rpc_port"`                // rpc port,default:9001
	WebApiPort          string `json:"web_api_port"`            // web api port,default:9002
	DataDir             string `json:"data_dir"`                // database local path,default:debrain-data/
	IsOpenSyncService   bool   `json:"is_open_sync_service"`    // is open sync service,default:true
	IsOpenRpcService    bool   `json:"is_open_sync_service"`    // is open rpc service,default:false
	IsOpenWebApiService bool   `json:"is_open_web_api_service"` // is open web api service,default:false
	IsOpenConsole       bool   `json:"is_open_console"`         // is open console,defaule false
}

func (op *Options) Start() {
	flag.Usage = usage
	args := os.Args
	for _, arg := range args {
		switch arg {
		case "console":
			{
				op.IsOpenConsole = true
			}
		}
	}
	flag.StringVar(&op.SyncPort, "syncport", "9000", "sync port")
	flag.StringVar(&op.RpcPort, "rpcport", "9001", "rpc port")
	flag.StringVar(&op.WebApiPort, "webapiport", "9002", "web api port")
	flag.StringVar(&op.DataDir, "dbdir", "debrain-data/", "sync port (necessary)")
	flag.BoolVar(&op.IsOpenSyncService, "opensync", true, "Whether to open the sync service")
	flag.BoolVar(&op.IsOpenRpcService, "openrpc", false, "Whether to open the rpc service")
	flag.BoolVar(&op.IsOpenWebApiService, "openwebapi", false, "Whether to open the webapi service")
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stdout, "debrain version: debrain/0.0.1 Usage: \n"+
		"debrain [-help] [-syncport port] [-rpcport port] [-webapiport port]\n"+
		"[-dbdir] [-opensync bool] [-openrpc bool] [-openwebapi bool]\n"+
		"Order: \n"+
		"   console\n"+
		"        Start an interactive command line\n"+
		"Options: \n")
	flag.PrintDefaults()
}
