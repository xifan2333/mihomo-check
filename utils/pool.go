package utils

import (
	"fmt"
	"sync"
)

type Task struct {
	ID     int
	Args   interface{}
	Result interface{}
	Err    error
}

type ThreadPool struct {
	NumWorkers int
	Tasks      chan Task
	Results    chan Task
	Wg         sync.WaitGroup
	TaskFunc   func(interface{}) (interface{}, error)
	results    []Task
	mu         sync.Mutex
	resultWg   sync.WaitGroup
	taskSentWg sync.WaitGroup
}

func NewThreadPool(numWorkers int, taskFunc func(interface{}) (interface{}, error)) *ThreadPool {
	tp := &ThreadPool{
		NumWorkers: numWorkers,
		Tasks:      make(chan Task, numWorkers*10),
		Results:    make(chan Task, numWorkers*10),
		TaskFunc:   taskFunc,
		results:    make([]Task, 0),
	}
	go tp.processResults()
	return tp
}

func (tp *ThreadPool) Start() {
	for i := 0; i < tp.NumWorkers; i++ {
		tp.Wg.Add(1)
		go tp.worker()
	}
}

func (tp *ThreadPool) worker() {
	defer tp.Wg.Done()
	for task := range tp.Tasks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					task.Err = fmt.Errorf("panic: %v", r)
					tp.Results <- task
				}
			}()
			task.Result, task.Err = tp.TaskFunc(task.Args)
			tp.Results <- task
		}()
	}
}

func (tp *ThreadPool) processResults() {
	for result := range tp.Results {
		tp.mu.Lock()
		tp.results = append(tp.results, result)
		tp.mu.Unlock()
		tp.resultWg.Done()
	}
}

func (tp *ThreadPool) GetResults() []Task {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return tp.results
}

func (tp *ThreadPool) AddTaskArgs(argsList []interface{}) {
	tp.resultWg.Add(len(argsList))
	tp.taskSentWg.Add(1)

	go func() {
		defer tp.taskSentWg.Done()
		for i, args := range argsList {
			task := Task{
				ID:   i + 1,
				Args: args,
			}
			tp.Tasks <- task
		}
	}()
}

func (tp *ThreadPool) Wait() {
	tp.taskSentWg.Wait()
	close(tp.Tasks)
	tp.Wg.Wait()
	close(tp.Results)
	tp.resultWg.Wait()
}
