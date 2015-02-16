package main

import (
	"fmt"
	"antigen/identity"
	"os"
	"bufio"
	
	)

func main(){
        scanner := bufio.NewScanner(os.Stdin)

        for scanner.Scan() {
	  GenAdd(scanner.Text())

	  
	}
	}


func GenAdd(password string){
	var id, _ = identity.NewDeterministic(password, 1)
	id.CreateAddress(4,1)
	address, signingkey, encryptionkey, _ := id.Export()
	fmt.Printf("{%q:{\"address\":%q,\"signingkey\":%q,\"encryptionkey\":%q}}\n", password, address, signingkey, encryptionkey)
	}

