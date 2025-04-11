package polling

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

// Process 将数据分成多个协程进行处理
func Process[T any](goRoutineSize int, data []T, f func(d T) error) {
	dataCh := make(chan T)
	defer close(dataCh)
	wg := &sync.WaitGroup{}
	wg.Add(len(data))
	executor := func(ctx context.Context) {
		for d := range dataCh {
			if err := f(d); err != nil {
				logrus.Errorf("[executor] - 执行失败, err = %v", err)
			}
			wg.Done()
		}
	}

	for i := 0; i < goRoutineSize; i++ {
		ctx := context.WithValue(context.Background(), "index", i)
		go executor(ctx)
	}

	for _, d := range data {
		dataCh <- d
	}
	wg.Wait()
}
