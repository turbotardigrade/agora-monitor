package main

import (
	"fmt"
	"time"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"

	"github.com/turbotardigrade/monitor/node"
)

func InitNode() (*core.IpfsNode, error) {
	// Need to increse limit for number of filedescriptors to
	// avoid running out of those due to a lot of sockets
	err := checkAndSetUlimit()
	if err != nil {
		return nil, err
	}

	// Create ipfs node if not exists
	addr := &config.Addresses{
		Swarm: []string{
			"/ip4/0.0.0.0/tcp/4004",
			"/ip6/::/tcp/4004",
		},
		API:     "/ip4/127.0.0.1/tcp/5004",
		Gateway: "/ip4/127.0.0.1/tcp/8084",
	}

	if !Exists(MyNodePath) {
		err := node.NewNodeRepo(MyNodePath, addr)
		if err != nil {
			return nil, err
		}

		fmt.Println("Seeding...")
		time.Sleep(5 * time.Second)
		fmt.Println("Seeding done.")

	}

	return node.NewNode(MyNodePath)
}
