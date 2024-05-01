package core

import (
	"fmt"
	"sync"
)

type fetchJob struct {
	key  string
	file UploadableFile
}

type fetchQueue chan *fetchJob

func SelectUploadFiles(workerNum int, que fetchQueue) chan UploadableFile {
	var wg sync.WaitGroup
	out := make(chan UploadableFile)

	for range workerNum {
		wg.Add(1)
		go fetch(&wg, que)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func fetch(wg *sync.WaitGroup, que fetchQueue) {
	defer wg.Done()
	for {
		job, ok := <-que
		if !ok {
			return
		}
		fmt.Println(job)
	}
}
