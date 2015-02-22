package main

import (

	"fmt"
	"os"
	"bufio"
	"runtime"
	"time"
	"sync"
	"sync/atomic"
	"flag"
	"errors"
	"compress/bzip2"

	"github.com/syndtr/goleveldb/leveldb"
	//"github.com/syndtr/goleveldb/leveldb/errors"
	//"github.com/syndtr/goleveldb/leveldb/opt"
	//"github.com/syndtr/goleveldb/leveldb/storage"
	//"github.com/syndtr/goleveldb/leveldb/table"
	//"github.com/syndtr/goleveldb/leveldb/util"
	//"log"


	"github.com/antitree/antigen/identity"
	
)

var checked uint64 = 0

var debug bool
var cpus int 
var ct int
var file string
var start uint64

func init() {

	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.IntVar(&cpus, "cpus", runtime.NumCPU(), "number of cpu threads")
	flag.IntVar(&ct, "ct", runtime.NumCPU()-1, "number of crypto threads")
	flag.StringVar(&file, "file", "", "file to read from, use '-' for stdin")

	flag.Parse()

	if debug {
		fmt.Printf("CPU:%d:CT:%d:FILE:%q\n", cpus, ct, file)
	}

	runtime.GOMAXPROCS(cpus)

}

type result struct {

	password string
	address string
	signingkey string
	encryptionkey string
	err error

}

func main(){

	start = uint64(time.Now().UnixNano())

	done := make(chan struct{})
	defer close(done)


	balls, errc := parseInput(done)
	if err := <-errc; err != nil {
		panic(err)
	}

	c := make(chan result, 1000)

	var wg sync.WaitGroup
	wg.Add(ct)

	for i := 0; i < ct ; i++ {

		if debug {
			fmt.Printf("Adding worker:%d\n", i)
		}

		go func() {

			worker(done, balls, c)
			wg.Done()

		}()

	}

	go func() {
		wg.Wait()
		close(c)
	}()

	db, err := leveldb.OpenFile("BallZ.db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()


	for r := range c  {
		fmt.Printf("{%q:{\"address\":%q,\"signingkey\":%q,\"encryptionkey\":%q}}\n", r.password, r.address, r.signingkey, r.encryptionkey)
		err = db.Put([]byte(r.address), []byte(r.password), nil)
		if err != nil {
			fmt.Printf("Error on put:%q\n", err)
		}
	}

	var stop uint64 = uint64(time.Now().UnixNano())

	tt := float64(stop-start)/1e9
	tt = float64(checked) / tt
	fmt.Printf("Processed %d words in %.3f seconds for a rate of %.4f/second\n", checked, float64(stop-start)/1e9, tt)

}

// digester reads path names from paths and sends digests of the corresponding
// files on c until either paths or done is closed.
func worker(done <-chan struct{}, balls <-chan string, c chan<- result) {

	for password := range balls { 

		var id, _ = identity.NewDeterministic(password, 1)
		id.CreateAddress(4,1)
		address, signingkey, encryptionkey, err := id.Export()

		atomic.AddUint64(&checked, 1)
   		var total uint64 = uint64(time.Now().UnixNano())

		var diff uint64 = (total - start)

		if debug == true {
			fmt.Printf("time:%.8f \n", float64( diff / checked ) /1e9)
		}

		select {
			case c <- result{password, address, signingkey, encryptionkey, err}:
			case <-done:
			return
		}
	}

}


//
//
//
func parseInput(done <-chan struct{} ) (<-chan string, <-chan error) {

	balls := make(chan string, 100)
	errc := make(chan error, 1)

	go func() { 

		// close when done
		defer close(balls)

		scanner := bufio.NewScanner(os.Stdin)
	    scanner.Split(bufio.ScanLines)


		if ( file != "-" ) {

			if debug == true {
				fmt.Printf("file:%q\n", file)
			}
			f, err := os.Open(file)
			if (err != nil) {
				errc <- err
				return
			}
			errc <- nil

			zReader := bzip2.NewReader(f)

			scanner = bufio.NewScanner(zReader)

		}

		for scanner.Scan() {

			password := scanner.Text()

			select {
				case balls <- password:
				case <-done:
					errc <- errors.New("cancelled")
				return
			}


		}


	}()

	return balls, errc

}







