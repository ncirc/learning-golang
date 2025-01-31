package main

import (
	"context"
	"log"
	"os"
	"time"
)

func scanDir(dir string, ch chan<- string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ch <- f.Name()
	}
}

func main() {
	chFile := make(chan string)
	finished := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	log.Println("starting scan")
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				finished <- struct{}{}
				return
			default:
				scanDir("/home/sascha/.temp/in", chFile)
				time.Sleep(15 * time.Second)
			}
		}
	}(ctx)

	go func() {
		for f := range chFile {
			log.Println(f)
		}
	}()

	go func() {
		time.Sleep(3 * time.Minute)
		cancel()
	}()

	<-finished
	log.Println("done")
}
