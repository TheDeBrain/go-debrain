package sync

import (
	"github.com/derain/core/protocols"
	"github.com/derain/internal/pkg/rules"
)


// sync system db table
func SynSysDBT() {
	p := new(protocols.CommProtocol)
	p.ProtocolType = rules.SYS_DB_SYNC_PROTOCOL
}

// sync file system db table
func SyncFileSysDBT() {
	p := new(protocols.CommProtocol)
	p.ProtocolType = rules.SYS_DB_SYNC_PROTOCOL
}
