package polling_test

import (
	"github.com/sirupsen/logrus"
	"otelDemo/utils/polling"
	"testing"
)

func TestProcess(t *testing.T) {
	f := func(d int) error {
		logrus.Infof("[f] - 处理数据, d = %d", d)
		return nil
	}

	polling.Process(5, []int{1, 2, 3, 4, 5}, f)
}
