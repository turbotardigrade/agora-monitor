package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/turbotardigrade/monitor/node"
)

const MyNodePath = "./data/monitorNode"

// Exists check if path exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func main() {
	// Need to increse limit for number of filedescriptors to
	// avoid running out of those due to a lot of sockets
	err := checkAndSetUlimit()
	if err != nil {
		panic(err)
	}

	if !Exists(MyNodePath) {
		err := node.NewNodeRepo(MyNodePath, nil)
		if err != nil {
			panic(err)
		}
	}

	n, err := node.NewNode(MyNodePath)
	if err != nil {
		panic(err)
	}

	res, err := node.Request(n, "QmdtfJBMitotUWBX5YZ6rYeaYRFu6zfXXMZP6fygEWK2iu", "/health", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

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
