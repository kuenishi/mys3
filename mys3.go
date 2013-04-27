
package main

import ("fmt"
	"runtime"
	"log"
)

func die(ch chan int) {
	defer func() {
		var err = recover()
		ch <- -127
		log.Println("work failed:", err)
		close(ch)
	}()
	panic("oops!\n")
	close(ch)
	ch <- 23
}

func hello(a, b, c int, ch chan int) {
	defer func() {
		var err = recover()
		ch <- -127
		log.Println(" >>> work failed:", err)
		close(ch)
	}()
	fmt.Printf("a %d b %d c %d\n", a, b, c)
	close(ch)
}

func add(x, y int) int {
	return x+y
}
func pa(f func(x, y int) int, i int) func(x int) int {
	g := func(j int) int {
		return f(i,j)
	}
	return g
}

func main() {
	fmt.Printf("pa -> %d\n", (pa(add, 1))(2))

	runtime.GOMAXPROCS(5)
	var ch = make(chan int, 10)
	var i int
	i = 2
	go hello(1,i,3, ch)
	//go die(ch)
	fmt.Printf("hello, world\n")
	i = <-ch
	fmt.Printf("recvd: %d\n", i)
}
