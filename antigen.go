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

	runtime.GOMAXPROCS(cpus)

	// read word list into a channel
	go func() {

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

		for scanner.Scan() {

			password := scanner.Text()
			if debug == true {
				fmt.Printf("new:%s\n", password)
			}

			balls <- password
		}

		fmt.Printf("Processed:%d\n", checked)
		sem <- 1

	}()


	for i := 0; i < ct ; i++ {

		go func() {

			var password string
			ok := true
			for {

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

				var id, _ = identity.NewDeterministic(password, 1)
				id.CreateAddress(4,1)
				address, signingkey, encryptionkey, _ := id.Export()
				fmt.Printf("{%q:{\"address\":%q,\"signingkey\":%q,\"encryptionkey\":%q}}\n", password, address, signingkey, encryptionkey)

				atomic.AddUint64(&checked, 1)
   				var total uint64 = uint64(time.Now().UnixNano())

				var diff uint64 = (total - start)

				if debug == true {
					fmt.Printf("time:%.8f \n", float64( diff / checked ) /1e9)
				}

			}

		}()

	}

	// stop here when we are done
	<- sem
	done = true

	var stop uint64 = uint64(time.Now().UnixNano())

	fmt.Printf("Total Time:%.3f\n", float64(stop-start)/1e9 )




}











