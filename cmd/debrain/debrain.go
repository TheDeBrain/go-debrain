package main

import (
	"github.com/derain/core/initializa"
	"github.com/derain/core/sync"
	"github.com/derain/core/task"
	"github.com/derain/internal/app/api/web"
)

func init() {
	//// init options
	//options := new(command.Options)
	//options.Start()
	// init db
	initializa.InitDB()
	task.Start()
	//// start sync service
	//if options.IsOpenSyncService {
	go sync.StartSyncService()
	//}
	//// start rpc service
	//if options.IsOpenRpcService {
	//	go rpc.StartRpcService()
	//}
	//// start web api service
	//if options.IsOpenWebApiService {
	web.StartWebApiService()
	//}
	//// init console
	//if options.IsOpenConsole {
	//	console.Start()
	//}
}

func main() {

}
