package multi_threading

import (
	"math/rand"
	"sync"
)

type mii map[int]int

type NetworkLayer struct {
	QueueMap  []*MessageQueue
	ErrorRate float64
	Rng       rand.Rand
	Tick      int64 // TODO: implemet volatile (sync.Atomic ??)
	StopFlag  bool
	//TODO: uniform_real_distribution<> Distrib = uniform_real_distribution<>(0.0, 1.0);

	networkSize int
	networkMap  map[int]mii
	//TODO(hard): globalTimer thread := thread(globalTimerExecutor, this);
	globalTimerMutex sync.Mutex // TODO(hard): implement recursuveness look https://habr.com/ru/post/271789/ (5 section) or RWMutex

}
