package util

import (
	"fmt"
	"log"
	"time"
)

func Retry(attempts int, sleepDuration time.Duration, f func() (error, bool)) error {
	for i := 0; i < attempts; i++ {
		err, done := f()
		if err != nil {
			return err
		}
		if done {
			return nil
		} else {
			log.Printf("[INFO] operation not finished yet, retrying in %v seconds\n", sleepDuration.Seconds())
			time.Sleep(sleepDuration)
		}
	}
	return fmt.Errorf("giving up after %v attempts", attempts)
}
