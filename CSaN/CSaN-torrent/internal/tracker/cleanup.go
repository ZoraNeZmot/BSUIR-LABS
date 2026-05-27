package tracker

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func StartCleanup(ctx context.Context, log *logrus.Logger, store *Store, timeout time.Duration) {
	if log == nil {
		log = logrus.StandardLogger()
	}
	ticker := time.NewTicker(timeout / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cutoff := time.Now().UTC().Add(-timeout)
			removed := store.RemoveStale(cutoff)
			if removed > 0 {
				log.WithFields(logrus.Fields{
					"removed": removed,
					"cutoff":  cutoff,
				}).Info("peer cleanup")
			}
		}
	}
}
