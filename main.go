package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"

	"github.com/turbotardigrade/monitor/node"
)

const MyNodePath = "./data/monitorNode"

var NodeList = []string{
	"QmdtfJBMitotUWBX5YZ6rYeaYRFu6zfXXMZP6fygEWK2iu",
	"QmVmPkKN9XXfxwQfinSWDYuU8M6U9dZdL46uSoSwuYgLKL",
}

func main() {
	// Need to increse limit for number of filedescriptors to
	// avoid running out of those due to a lot of sockets
	err := checkAndSetUlimit()
	if err != nil {
		panic(err)
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
			panic(err)
		}

		fmt.Println("Seeding...")
		time.Sleep(5 * time.Second)
		fmt.Println("Seeding done.")

	}

	n, err := node.NewNode(MyNodePath)
	if err != nil {
		panic(err)
	}

	healthy := make(map[string]bool, len(NodeList))
	posts := make(map[string][]string, len(NodeList))

	var wg sync.WaitGroup
	wg.Add(len(NodeList))
	for _, target := range NodeList {
		go func(target string) {
			defer wg.Done()

			ps := getPosts(n, target)
			if ps != nil {
				posts[target] = ps
				healthy[target] = true
			} else {
				healthy[target] = false
			}
		}(target)
	}

	wg.Wait()

	fmt.Println("Health Status")
	for k, v := range healthy {
		fmt.Println(k, v)
	}

	fmt.Println("\nPosts")
	for k, v := range posts {
		fmt.Println(k, v)
	}

}

func getPosts(n *core.IpfsNode, target string) []string {
	resp, err := node.Request(n, target, "/posts")
	if err != nil {
		fmt.Println("Request failed:", err)
		return nil
	}

	js := make(map[string][]string)
	err = json.Unmarshal([]byte(resp), &js)
	if err != nil {
		fmt.Println("JSON unmarshalling failed:", err)
		return nil
	}

	return js["Posts"]
}
