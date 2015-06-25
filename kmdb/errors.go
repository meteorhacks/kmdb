package kmdb

func (c ErrorCode) Error() (why string) {
	switch c {
	case ERR_NO_ERROR:
		return "ERR_NO_ERROR"
	case ERR_DB_NOT_FOUND:
		return "ERR_DB_NOT_FOUND"
	case ERR_BATCH_ERROR:
		return "ERR_BATCH_ERROR"
	default:
		return "ERR_UNKNOWN"
	}
}
