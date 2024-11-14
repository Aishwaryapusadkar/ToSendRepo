package utils

import (
	"sync"
)

var (
	sequenceNum = uint32(0)
	lock        sync.Mutex
)

func GetSequenceNumber() uint32 {
	lock.Lock()
	defer lock.Unlock()
	sequenceNum++
	return sequenceNum
}
