package worker

import (
	"math/rand"
	"time"
)

type SumJob struct {
	Numbers []int
}

func (sumJob SumJob) Execute() interface{} {
	sum := 0
	for _, number := range sumJob.Numbers {
		sum += number
	}

	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	return sum
}
