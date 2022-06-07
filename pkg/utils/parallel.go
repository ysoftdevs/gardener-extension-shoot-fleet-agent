package utils

import (
	"github.com/pkg/errors"
	"sync"
)

// RunParallelFunctions runs a list of functions in parallel
// returns nil if all functions return nil
// returns an error which wraps all errors that occurred within the functions
func RunParallelFunctions(functions []func() error) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorList []error

	for _, function := range functions {
		wg.Add(1)

		// create an intermediate variable so each goroutine uses the correct function and doesn't get overwritten
		function := function
		go func(wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			err := function()
			if err != nil {
				mu.Lock()
				errorList = append(errorList, err)
				mu.Unlock()
			}
		}(&wg, &mu)
	}

	wg.Wait()

	if len(errorList) > 0 {
		err := errors.New("not all functions returned nil error")
		for _, errorItem := range errorList {
			err = errors.Wrap(err, errorItem.Error())
		}
		return err
	}

	return nil
}
