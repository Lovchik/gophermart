package retry

import (
	log "github.com/sirupsen/logrus"
	"time"
)

var AttemptsDelay = 2 * time.Second

func Retry[T any](attempts int, delay time.Duration, function func() (T, error)) (T, error) {
	var err error
	var result T
	for i := 0; i < attempts; i++ {
		log.Info("Retrying ", i+1)
		result, err = function()
		if err == nil {
			if attempts == attempts-1 {
				log.Info("Successfully ran ", result)
			}
			return result, nil
		} else {
			log.Error(err)
		}

		if i < attempts-1 {
			time.Sleep(delay)
		}
		delay += AttemptsDelay
	}

	return result, err
}
