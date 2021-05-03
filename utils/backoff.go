package utils

import (
	"github.com/apex/log"
	"time"
)

func BackoffErrorNotify(err error, duration time.Duration) {
	log.Errorf("backoff error: %v, to retry in %1f seconds", err, duration.Seconds())
}
