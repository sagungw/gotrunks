package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-redsync/redsync/v4"
)

var (
	ResourceKeyHeader = "X-Lock-Resource-Key"
)

type MutexMiddlewareConfig struct {
	Retry  int `mapstructure:"retry"`
	TtlSec int `mapstructure:"ttlsec"`
}

func NewMutexMiddleware(config *MutexMiddlewareConfig, rs *redsync.Redsync) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resourceKey := r.Header.Get(ResourceKeyHeader)
			mutex := rs.NewMutex(
				fmt.Sprintf("lock:generic:%s", resourceKey),
				redsync.WithExpiry(time.Duration(config.TtlSec)*time.Second),
				redsync.WithTries(config.Retry),
			)

			if err := mutex.LockContext(r.Context()); err != nil {
				slog.ErrorContext(r.Context(), "failed acquiring lock", "resource", resourceKey, "error", err)
				http.Error(w, fmt.Sprintf("failed acquiring lock for resource %s", resourceKey), http.StatusLocked)
				return
			}

			defer func() {
				if ok, err := mutex.UnlockContext(r.Context()); !ok || err != nil {
					slog.ErrorContext(r.Context(), "failed releasing lock", "resource", resourceKey, "error", err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
