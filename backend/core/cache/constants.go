package cache

import "time"

const (
	TTLDefault    = 5 * time.Minute
	TTLShort      = 1 * time.Minute
	TTLLong       = 30 * time.Minute
	TTLSession    = 7 * 24 * time.Hour

	KeySeparator  = ":"
)
