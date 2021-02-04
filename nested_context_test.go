package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestNestedContext(t *testing.T){
	var (
		wg sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subCtx := context.WithValue(ctx, "name", "index")
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(ctx context.Context, n int) {
			defer wg.Done()
			name, bOk := ctx.Value("name").(string)
			if !bOk {
				fmt.Println("type of context's value is not int")
				return
			}
			if 0 == n {
				myCtx, myCancel := context.WithCancel(subCtx)
				go func() {
					select {
					case <-myCtx.Done():
						fmt.Println("---------------- close")
						return
					}
				}()
				t := time.NewTimer(time.Second * 3)
				select {
				case <-t.C:
					myCancel()
				}
			}
			frame := time.NewTicker(time.Second * 1)
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("name %s_%d exit\n", name, n)
					return
				case <-frame.C:
					fmt.Printf("name %s_%d is running\n", name, n)
				}
			}
		}(subCtx, i)
	}
	waitForASignal1()
	cancel()
	wg.Wait()
	fmt.Println("context test over")
}

func waitForASignal1()  {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
}