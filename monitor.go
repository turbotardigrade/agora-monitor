package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"

	"github.com/turbotardigrade/monitor/node"
)

func monitor(n *core.IpfsNode) (healthy map[string]bool, posts map[string][]string) {
	healthy = make(map[string]bool, len(NodeList))
	posts = make(map[string][]string, len(NodeList))

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

	return healthy, posts
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
