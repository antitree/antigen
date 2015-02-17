package main

import (

	"fmt"
	"github.com/antitree/antigen/identity"
	"os"
	"bufio"
	"runtime"
	"time"
	"sync/atomic"
	"flag"
	
)


var balls = make(chan string, 100)
var sem = make(chan int)
var done = false

var checked uint64 = 0

var debug = flag.Bool("debug", false, "enable debug logging")
var cpus = flag.Int("cpu", runtime.NumCPU(), "number of cpu threads")
var ct = flag.Int("ct", runtime.NumCPU()-1, "number of crypto threads")

func main(){

	flag.Parse()

	var start uint64 = uint64(time.Now().UnixNano())

	runtime.GOMAXPROCS(*cpus)

	f, err := os.Open("crackstation.txt")
	if (err != nil) {
		panic(err)
	}

	// read word list into a channel
	go func() {

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {

			password := scanner.Text()
			if *debug == true {
				fmt.Printf("new:%s\n", password)
			}

			balls <- password
		}

		fmt.Printf("Doneski!\n")
		sem <- 1

	}()


	for i := 0; i < *ct ; i++ {

		go func() {

			for {

				if done == true { break ; }

				password := <- balls
				var id, _ = identity.NewDeterministic(password, 1)
				id.CreateAddress(4,1)
				address, signingkey, encryptionkey, _ := id.Export()
				fmt.Printf("{%q:{\"address\":%q,\"signingkey\":%q,\"encryptionkey\":%q}}\n", password, address, signingkey, encryptionkey)

				atomic.AddUint64(&checked, 1)
   				var total uint64 = uint64(time.Now().UnixNano())

				var diff uint64 = (total - start)

				if *debug == true {
					fmt.Printf("time:%.8f \n", float64( diff / checked ) /1e9)
				}

			}

		}()

	}

	// stop here when we are done
	<- sem
	done = true

	var stop uint64 = uint64(time.Now().UnixNano())

	fmt.Printf("Total Time:%.3f", float64(stop-start)/1e9 )




}











