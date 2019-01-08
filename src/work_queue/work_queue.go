package work_queue

type Worker interface {
	Run() interface{}
}

type WorkQueue struct {
	Jobs    chan Worker
	Results chan interface{}
	Workers uint	// nWorkers
	StopRequests chan uint	// Stop request
}

// Create a new work queue capable of doing nWorkers simultaneous tasks, expecting to queue maxJobs tasks.
func Create(nWorkers uint, maxJobs uint) *WorkQueue {
	wq := new(WorkQueue)

	// initialize struct
	wq.Jobs = make(chan Worker, maxJobs)
	wq.Results = make(chan interface{}, maxJobs)
	wq.Workers = nWorkers
	wq.StopRequests = make(chan uint, maxJobs)

	// start nWorkers workers as goroutines
	for i := uint(0); i < nWorkers; i++ {
		go wq.worker()
	}

	return wq
}

// A worker goroutine that processes tasks from .Jobs unless .StopRequests has a message saying to halt now.
func (queue WorkQueue) worker() {
	// Listen on the .Jobs channel for incoming tasks.
	for {
		if queue.Jobs != nil && len(queue.StopRequests) == 0 {
			for i := range queue.Jobs {
				if len(queue.StopRequests) == 0 {
					// run tasks by calling .Run() & send the return value back on Results channel.
					queue.Results <- i.Run()
				}
			}
		} else {
			break
		}
	}

	return // Exit (return) when .Jobs is closed.
}

// put the work into the Jobs channel so a worker can find it and start the task.
func (queue WorkQueue) Enqueue(work Worker) {
	queue.Jobs <- work
}

// close .Jobs and remove all remaining jobs from the channel.
func (queue WorkQueue) Shutdown() {
	queue.StopRequests <- queue.Workers
}
