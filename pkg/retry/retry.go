// Package retry provides a function to execute an action with retries using specified delays between attempts.
package retry

import "time"

// Retry executes the given action repeatedly until success or exhaustion of retries.
// It waits for specified durations (`cooldowns`) between each failed attempt.
// Optionally calls `errCallback` on failure with the current error and retry count.
// Returns the final error encountered or nil if successful.
//
// The Retry function takes three arguments:
//   - action: A function that returns an error if the operation fails.
//   - cooldowns: A slice of time durations specifying how long to wait before each retry attempt.
//   - errCallback: An optional callback function invoked when an error occurs during execution.
func Retry(action func() error, cooldowns []time.Duration, errCallback func(error, uint8)) error {
	var err error
	for idx, delay := range cooldowns {
		err = action()
		if err == nil {
			return nil
		}
		if errCallback != nil {
			errCallback(err, uint8(idx+1))
		}
		time.Sleep(delay)
	}
	return action()
}
