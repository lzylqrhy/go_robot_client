package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNestedContext(t *testing.T){
	var (
		wg sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sCtx := context.WithValue(ctx, "rand", "random")
	subCtx := context.WithValue(sCtx, "name", "index")
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(ctx context.Context, n int) {
			defer wg.Done()
			rd, bOk := ctx.Value("rand").(string)
			if bOk {
				fmt.Println("rand key = ", rd)
			}
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
	//waitForASignal()
	time.Sleep(time.Second * 5)
	cancel()
	wg.Wait()
	fmt.Println("context test over")
}

func TestMapChannel(t *testing.T) {
	var (
		wg sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	ch := make(chan int, 3)
	go func() {
		defer wg.Done()
		for n := range ch {
			//select {
			//case <-ctx.Done():
			//	return
			//default:
				fmt.Println(n)
			//}
		}
		fmt.Println("It's over 1")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("It's over")
				return
			default:
				ch<- rand.Int()
				time.Sleep(2*time.Second)
			}
		}
		fmt.Println("It's over 2")
	}()
	time.Sleep(time.Second * 10)
	cancel()
	close(ch)
	wg.Wait()
	fmt.Println("channel test over")
}

type base struct {
	n int
	U int
}

func (p *base) SetN()  {
	p.n = 10
}

func (p *base) Do()  {
	p.SetN()
}

type derive struct {
	m int
	base
}

func (p *derive) SetN()  {
	//p.base.SetN()
	p.m = 100
}

func (p *derive) Print()  {
	fmt.Println(*p)
}

func TestInherit(t *testing.T) {
	var d derive
	d.U = 10
	d.Do()
	d.Print()
}

func TestOther(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i:=0; i<10; i++ {
		rn := rand.Uint32()
		fmt.Println(rn)
	}
	b := []byte("wo ai")
	n := b[:1+copy(b[1:], b[3:])]
	fmt.Println(n)

	m := make(map[int]int)
	if m[4] > 1 {
		m[4] = 2
	}
	fmt.Println(m[4])
	m[1] = 3
	m[2] = 4
	m[3] = 5
	for k, v := range m {
		fmt.Println("map element", v)
		if 2 == k {
			delete(m, k)
		}
	}
	fmt.Println(m)
}