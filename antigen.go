package main

import (

	"fmt"
	"antigen/identity"
	"os"
	"bufio"
	"runtime"
	"time"
	"sync/atomic"
<<<<<<< HEAD
	//"sync"
=======
	"flag"
>>>>>>> aab2c608c4f035f68deeb9122a588e0bb7ff40fa
	
)


var balls = make(chan string, 100)  // setup comms between routines
var sem = make(chan int)	    // sem?
var done = false

var checked uint64 = 0
//var wg sync.WaitGroup

var debug bool
var cpus int 
var ct int
var file string

func init() {

	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.IntVar(&cpus, "cpus", runtime.NumCPU(), "number of cpu threads")
	flag.IntVar(&ct, "ct", runtime.NumCPU()-1, "number of crypto threads")
	flag.StringVar(&file, "file", "", "file to read from, use '-' for stdin")

	flag.Parse()

}

func main(){


	var start uint64 = uint64(time.Now().UnixNano())

<<<<<<< HEAD
	//var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU())

	f, err := os.Open("shortlist.txt")
	if (err != nil) {
		panic(err)
	}
=======
	runtime.GOMAXPROCS(cpus)
>>>>>>> aab2c608c4f035f68deeb9122a588e0bb7ff40fa

	// read word list into a channel
	go func() {

<<<<<<< HEAD
		//scanner := bufio.NewScanner(os.Stdin)
		//_ = f
		scanner := bufio.NewScanner(f)
=======
		var scanner *bufio.Scanner

		if ( file == "-" ) {
			scanner = bufio.NewScanner(os.Stdin)
		} else {
			f, err := os.Open(file)
			if (err != nil) {
				panic(err)
			}
			scanner = bufio.NewScanner(f)
		}
>>>>>>> aab2c608c4f035f68deeb9122a588e0bb7ff40fa

		for scanner.Scan() {

			password := scanner.Text()
			if debug == true {
				fmt.Printf("new:%s\n", password)
			}

			balls <- password
		}

		fmt.Printf("Doneski!\n")
		sem <- 1

	}()


<<<<<<< HEAD
	for i := 0; i < runtime.NumCPU()-2 ; i++ {
		//wg.Add(1)
=======
	for i := 0; i < ct ; i++ {
>>>>>>> aab2c608c4f035f68deeb9122a588e0bb7ff40fa

		go func() {

			var password string
			ok := true
			for {
<<<<<<< HEAD
				if done == true { break ; }
				//wg.Add(1)
=======

				select {

					case password, ok = <-balls:
					if ok {
						// balls
					} else {
						if done == true { break; }
					}
					default:
						// no balls
				}
>>>>>>> aab2c608c4f035f68deeb9122a588e0bb7ff40fa

				var id, _ = identity.NewDeterministic(password, 1)
				id.CreateAddress(4,1)
				address, signingkey, encryptionkey, _ := id.Export()
				fmt.Printf("{%q:{\"address\":%q,\"signingkey\":%q,\"encryptionkey\":%q}}\n", password, address, signingkey, encryptionkey)

				atomic.AddUint64(&checked, 1)
   				var total uint64 = uint64(time.Now().UnixNano())

				var diff uint64 = (total - start)

<<<<<<< HEAD
				fmt.Printf("time:%.8f \n", float64( diff / checked ) /1e9)
				//wg.Done()
=======
				if debug == true {
					fmt.Printf("time:%.8f \n", float64( diff / checked ) /1e9)
				}
>>>>>>> aab2c608c4f035f68deeb9122a588e0bb7ff40fa

			}

		}()

	}

	// stop here when we are done
	<- sem
	done = true

	var stop uint64 = uint64(time.Now().UnixNano())

	fmt.Printf("Total Time:%.3f", float64(stop-start)/1e9 )




}











