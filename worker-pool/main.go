package main

// Task represents task
type Task struct{}

// Perform performs task
func (t *Task) Perform() {}

// generateTasks generates tasks
func generateTasks(count int) []*Task {
	var tt []*Task
	for i := 0; i < count; i++ {
		t := &Task{}
		tt = append(tt, t)
	}
	return tt
}

type token struct{}

func main() {
	// define task count and generate tasks
	taskCount := 100
	tt := generateTasks(taskCount)
	// define limit and create semaphore
	limit := 10
	semaphore := make(chan *token, limit)
	// perform tasks if semaphore has enough space
	for _, t := range tt {
		semaphore <- &token{}
		go func(t *Task) {
			t.Perform()
			<-semaphore
		}(t)
	}
	// wait until all tasks are completed
	for n := limit; n > 0; n-- {
		semaphore <- &token{}
	}
}
