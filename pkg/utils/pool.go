package utils

import (
	"bse-eti-stream/pkg/logger"
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var (
	wp *WorkerPool
)

type WorkerPool struct {
	context  context.Context
	Capacity int
	Pool     []chan interface{}
	wg       *sync.WaitGroup
}

type WorkerPoolInterface interface {
	EnQueue(context.Context, int, interface{})
	DeQueue(context.Context, int, func(interface{}) (interface{}, error))
	Run(context.Context, func(interface{}) (interface{}, error))
	Size() int
	Wait()
}

func CreateWorkerPool(c context.Context, size int) {
	if size == 0 {
		size = 5
	}
	wp = &WorkerPool{
		context:  c,
		Capacity: size,
		Pool:     make([]chan interface{}, size),
		wg:       &sync.WaitGroup{},
	}
	for i := 0; i < size; i++ {
		wp.Pool[i] = make(chan interface{}, size)
	}
	logger.Log(c).Info("worker pool is created of size", zap.Int("size", size))
}

func GetWorkerPool() WorkerPoolInterface {
	if wp == nil {
		fmt.Println("WorkerPool is not initialized. Call CreateWorkerPool() first.")
		return nil
	}
	return wp
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}

func (w *WorkerPool) Size() int {
	return w.Capacity
}

func (w *WorkerPool) Run(c context.Context, task func(interface{}) (interface{}, error)) {
	if w == nil {
		fmt.Println("WorkerPool instance is nil")
		return
	}
	fmt.Println("###############message pool")

	// w.wg.Add(w.Capacity)

	// for i := 0; i < w.Capacity; i++ {
	// 	logger.Log(c).Info("worker started", zap.Int("index", i))
	// 	go w.DeQueue(c, i, task)
	// }
	// go w.DeQueue(c, 0, task)

	// w.wg.Wait()
}

func (w *WorkerPool) EnQueue(c context.Context, index int, data interface{}) {
	if w == nil || w.Pool[index] == nil {
		logger.Log(c).Error("Attempted to enqueue on an uninitialized worker pool.")
		return
	}
	w.Pool[index] <- data
	logger.Log(c).Debug("Added to Queue", zap.Any("indexAt", index))
}

func (w *WorkerPool) DeQueue(c context.Context, index int, task func(interface{}) (interface{}, error)) {
	defer w.wg.Done()
	data := <-w.Pool[index]
	response, err := task(data)
	logger.Log(c).Debug("Fetched from Queue", zap.Any("id", index), zap.Any("data", data), zap.Any("response", response), zap.Error(err))
}
