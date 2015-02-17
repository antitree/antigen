package main

import (
	"fmt"
	"github.com/antitree/antigen/identity"
	"os"
	"bufio"
	"runtime"
	"time"
	
	)


var balls = make(chan string, 100)
var sem = make(chan int)
var done = false

func main(){

	start := time.Now().UnixNano()
   	total := time.Now().UnixNano()
	diff := total - start
	var checked int64 = 0

	runtime.GOMAXPROCS(runtime.NumCPU())

	f, err := os.Open("crackstation.txt")
	if (err != nil) {
		panic(err)
	}

	// read word list into a channel
	go func() {

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			password := scanner.Text()
			fmt.Printf("new:%s\n", password)
			balls <- password
		}

		fmt.Printf("Doneski!\n")
		sem <- 1

	}()


	for i := 0; i < runtime.NumCPU()-2 ; i++ {

		go func() {

			for {

				if done == true { break ; }

				password := <- balls
				var id, _ = identity.NewDeterministic(password, 1)
				id.CreateAddress(4,1)
				address, signingkey, encryptionkey, _ := id.Export()
				fmt.Printf("{%q:{\"address\":%q,\"signingkey\":%q,\"encryptionkey\":%q}}\n", password, address, signingkey, encryptionkey)

				checked++
			   	total = time.Now().UnixNano()

				diff = total - start
				fmt.Printf("time:%.8f \n", float64( diff / checked ) /1e9)

			}

		}()

	}

	// stop here when we are done
	<- sem
	done = true

	stop := time.Now().UnixNano()

	fmt.Printf("Total Time:%.3f", float64(stop-start)/1e9 )




}











