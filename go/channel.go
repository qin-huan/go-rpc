package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		ch = make(chan []string, 0)
		wg = &errgroup.Group{}
	)
	wg.Go(func() error {
		foreach:
		for {
			select {
			case strs, ok := <-ch:
				if !ok {
					break foreach
				}
				fmt.Println(strs)
			}
		}
		return nil
	})

	var strs = []string{
		"hello",
	}
	ch <- strs
	close(ch)

	if err := wg.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("done")
}