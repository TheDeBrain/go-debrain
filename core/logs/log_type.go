package logs

const (
	SYNC_LOG = iota
)

var (
	LOG_TYPE_MAPPING = map[int]string{
		SYNC_LOG:  "sync_log.json",
	}
)
