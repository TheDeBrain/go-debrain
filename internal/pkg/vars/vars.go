package vars

import "github.com/derain/core/db/table/sys"

var (
	TSys = sys.LoadTSys()
	TFSys = sys.LoadFileSys()
)

