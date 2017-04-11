package main

import (
	"fmt"
	"log"
	"syscall"
)

var ipfsFileDescNum = uint64(5120)

// Taken from github.com/ipfs/go-ipfs/blob/master/cmd/ipfs/ulimit_unix.go
func checkAndSetUlimit() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error getting rlimit: %s", err)
	}

	if rLimit.Cur < ipfsFileDescNum {
		if rLimit.Max < ipfsFileDescNum {
			log.Println("Error: adjusting max")
			rLimit.Max = ipfsFileDescNum
		}
		// Info.Println("Adjusting current ulimit to ", ipfsFileDescNum, "...")
		rLimit.Cur = ipfsFileDescNum
	}

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return fmt.Errorf("Error setting ulimit: %s", err)
	}

	return nil
}
