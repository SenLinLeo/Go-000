package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 5; i++ {
		index := i
		group.Go(func() error {

			fmt.Printf("start to execute the %d gorouting\n", index)
			time.Sleep(time.Duration(index) * time.Second)
			if index%2 == 0 {
				return fmt.Errorf("something has failed on grouting:%d", index)
			}
			fmt.Printf("gorouting:%d end\n", index)
			return nil

		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}
