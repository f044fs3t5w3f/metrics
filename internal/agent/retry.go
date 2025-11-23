package agent

import "time"

func retry(action func() error, cooldowns []time.Duration, errCallback func(error, uint8)) error {
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
