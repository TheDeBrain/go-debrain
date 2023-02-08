package command

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	SyncPort            string `json:"sync_port"`               // sync port,default:9000
	RpcPort             string `json:"rpc_port"`                // rpc port,default:9001
	WebApiPort          string `json:"web_api_port"`            // web api port,default:9002
	DataDir             string `json:"data_dir"`                // database local path,default:debrain-data/
	IsOpenSyncService   bool   `json:"is_open_sync_service"`    // is open sync service,default:open
	IsOpenRpcService    bool   `json:"is_open_sync_service"`    // is open rpc service,default:close
	IsOpenWebApiService bool   `json:"is_open_web_api_service"` // is open web api service,default:close
}

func (c *Command) Start() {
	flag.Usage = usage
	flag.StringVar(&c.SyncPort, "syncport", "9000", "sync port")
	flag.StringVar(&c.RpcPort, "rpcport", "9001", "rpc port")
	flag.StringVar(&c.WebApiPort, "webapiport", "9002", "web api port")
	flag.StringVar(&c.DataDir, "dbdir", "debrain-data/", "sync port (necessary)")
	flag.BoolVar(&c.IsOpenSyncService, "opensync", true, "Whether to open the sync service")
	flag.BoolVar(&c.IsOpenRpcService, "openrpc", false, "Whether to open the rpc service")
	flag.BoolVar(&c.IsOpenWebApiService, "openwebapi", false, "Whether to open the webapi service")
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stdout, "debrain version: debrain/0.0.1 Usage: \n"+
		"debrain [-help] [-syncport port] [-rpcport port] [-webapiport port]\n"+
		"[-dbdir] [-opensync bool] [-openrpc bool] [-openwebapi bool]\n"+
		"Options: \n")
	flag.PrintDefaults()
}
