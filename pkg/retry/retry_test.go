package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestFunc(okAttempt int) func() error {
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

func TestRetry(t *testing.T) {
	t.Run("First attempt", func(t *testing.T) {
		testFunc := getTestFunc(1)
		now := time.Now()
		err := Retry(testFunc, []time.Duration{1 * time.Second}, nil)
		duration := time.Since(now)
		assert.NoError(t, err)
		seconds := duration.Seconds()
		assert.Less(t, seconds, 0.9)
	})
	t.Run("Second attempt", func(t *testing.T) {
		testFunc := getTestFunc(2)
		now := time.Now()
		err := Retry(testFunc, []time.Duration{1 * time.Second}, nil)
		duration := time.Since(now)
		assert.NoError(t, err)
		seconds := duration.Seconds()
		assert.Less(t, seconds, 1.9)
		assert.Less(t, 0.9, seconds)
	})
	t.Run("No success", func(t *testing.T) {
		testFunc := getTestFunc(3)
		err := Retry(testFunc, []time.Duration{1 * time.Second}, nil)
		assert.Error(t, err)
	})

}

func ExampleRetry_success() {
	f := getTestFunc(2)
	err := Retry(f, []time.Duration{1 * time.Second}, nil)
	fmt.Println(err)
}

func ExampleRetry_fail() {
	f := getTestFunc(3)
	err := Retry(f, []time.Duration{1 * time.Second}, nil)
	fmt.Println(err)
}
