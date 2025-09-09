package time

import (
	"log/slog"
	"time"
)

var (
	LocationJakarta *time.Location

	DefaultTimeNowUTCF = func() time.Time {
		return time.Now().UTC()
	}

	TimeNowUTCF = DefaultTimeNowUTCF
)

func TimeNowUTC() time.Time {
	return TimeNowUTCF()
}

func init() {
	l, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Error("failed to load location", "error", err, "location", "Asia/Jakarta")
	}

	LocationJakarta = l
}
