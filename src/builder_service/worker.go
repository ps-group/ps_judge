package main

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	workerNum     = 3
	sleepInterval = 1 * time.Second
)

type Task interface {
	Run(workerID int) error
}

type TaskGenerator interface {
	Next() (bool, Task)
}

func TaskProvider(generator TaskGenerator, stopChan chan struct{}) <-chan Task {
	tasksChan := make(chan Task)
	go func() {
		active := true
		for active {
			select {
			case <-stopChan:
				close(tasksChan)
				active = false
			default:
				ok, task := generator.Next()
				if ok {
					tasksChan <- task
				} else {
					time.Sleep(sleepInterval)
				}
			}
		}
	}()
	return tasksChan
}

func RunTaskProvider(generator TaskGenerator, stopChan chan struct{}) <-chan Task {
	resultChan := make(chan Task)
	stopTaskProviderChan := make(chan struct{})
	taskProviderChan := TaskProvider(generator, stopTaskProviderChan)
	onStop := func() {
		stopTaskProviderChan <- struct{}{}
		close(resultChan)
	}

	go func() {
		for {
			select {
			case <-stopChan:
				onStop()
				return
			case task := <-taskProviderChan:
				select {
				case <-stopChan:
					onStop()
					return
				case resultChan <- task:
				}
			}
		}
	}()
	return resultChan
}

func Worker(taskChan <-chan Task, workerID int) {
	for task := range taskChan {
		err := task.Run(workerID)
		if err != nil {
			logrus.WithField("error", err).Error("failed to execute task")
		}
	}
}

func RunWorkerPool(generator TaskGenerator, stopChan chan struct{}) *sync.WaitGroup {
	var wg sync.WaitGroup
	tasksChan := RunTaskProvider(generator, stopChan)
	for i := 0; i < workerNum; i++ {
		go func(i int) {
			wg.Add(1)
			Worker(tasksChan, i)
			wg.Done()
		}(i)
	}
	return &wg
}
