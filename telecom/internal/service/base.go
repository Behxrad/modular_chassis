package service

import (
	"sync"
)

type Base struct{}

var (
	once     sync.Once
	implsIns *Base
)

func GetBaseService() *Base {
	once.Do(func() {
		if implsIns == nil {
			implsIns = &Base{}
		}
	})
	return implsIns
}
