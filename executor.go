package main

func ExecuteInParallel(tasks []func() error, limit int, statusHandler func(complete, total int)) error {
	if limit == 0 {
		limit = len(tasks)
	}
	if limit > len(tasks) {
		limit = len(tasks)
	}
	wait := make(chan error, limit)
	running := 0
	complete := 0
	for _, task := range tasks {
		go func(t func() error) {
			wait <- t()
		}(task)
		running++
		if running == limit {
			err := <-wait
			running--
			complete++
			statusHandler(complete, len(tasks))
			if err != nil {
				return err
			}
		}
	}
	for running > 0 {
		err := <-wait
		running--
		complete++
		statusHandler(complete, len(tasks))
		if err != nil {
			return err
		}
	}
	return nil
}
