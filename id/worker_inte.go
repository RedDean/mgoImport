package id

import "sync"

type Worker interface {
	Do(<-chan interface{}, *sync.WaitGroup)
}
