package initialize

import (
	"sync"
)

func Run() {
	list := []func(){
		runHttpSrver,
	}

	var wg sync.WaitGroup

	for _, f := range list {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(f)
	}

	wg.Wait()
}
