package agent

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	getTestFunc := func(okAttempt int) func() error {
		var attempt int
		testFunc := func() error {
			attempt += 1
			if attempt == okAttempt {
				return nil
			}
			return errors.New("Ошибочка")
		}
		return testFunc
	}

	t.Run("First attempt", func(t *testing.T) {
		testFunc := getTestFunc(1)
		now := time.Now()
		err := retry(testFunc, []time.Duration{1 * time.Second}, nil)
		duration := time.Since(now)
		assert.NoError(t, err)
		seconds := duration.Seconds()
		assert.Less(t, seconds, 0.9)
	})
	t.Run("Second attempt", func(t *testing.T) {
		testFunc := getTestFunc(2)
		now := time.Now()
		err := retry(testFunc, []time.Duration{1 * time.Second}, nil)
		duration := time.Since(now)
		assert.NoError(t, err)
		seconds := duration.Seconds()
		assert.Less(t, seconds, 1.9)
		assert.Less(t, 0.9, seconds)
	})
	t.Run("No success", func(t *testing.T) {
		testFunc := getTestFunc(3)
		err := retry(testFunc, []time.Duration{1 * time.Second}, nil)
		assert.Error(t, err)
	})

}
