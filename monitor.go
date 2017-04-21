package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
)

var MonitorMap = make(map[string]MonitorResp)
var MonitorMapLock sync.RWMutex

func monitor(n *core.IpfsNode) (healthy map[string]bool, posts map[string][]string, blacklists map[string][]map[string]float32) {
	healthy = make(map[string]bool, len(NodeList))
	posts = make(map[string][]string, len(NodeList))
	blacklists = make(map[string][]map[string]float32, len(NodeList))

	var pslock sync.Mutex
	var bllock sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(NodeList))
	for _, target := range NodeList {
		go func(target string) {
			defer wg.Done()

			ps := getPosts(n, target)

			pslock.Lock()
			if ps != nil {
				posts[target] = ps
				healthy[target] = true
			} else {
				healthy[target] = false
			}
			pslock.Unlock()

			bl := getBlacklist(n, target)
			if bl != nil {
				bllock.Lock()
				blacklists[target] = bl
				bllock.Unlock()
			}
		}(target.ID)
	}
	wg.Wait()

	return healthy, posts, blacklists
}

func getPosts(n *core.IpfsNode, target string) []string {
	resp, err := Request(n, target, "/posts")
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

func getBlacklist(n *core.IpfsNode, target string) []map[string]float32 {
	resp, err := Request(n, target, "/blacklist")
	if err != nil {
		fmt.Println("Request failed:", err)
		return nil
	}

	js := make(map[string][]map[string]float32)
	err = json.Unmarshal([]byte(resp), &js)
	if err != nil {
		fmt.Println("JSON unmarshalling failed:", err)
		return nil
	}

	return js["Peers"]
}

func evalWorker(n *core.IpfsNode, hashes chan string, labels chan<- bool) {
	for hash := range hashes {
		content, err := getContent(n, hash)
		if err != nil {
			fmt.Println("ERROR getting content:", err)
			continue
		}

		label, err := checkLabel(content)
		if err != nil {
			fmt.Println("ERROR getting label:", err)
			continue
		}

		labels <- label
	}
}

func evaluatePosts(n *core.IpfsNode, posts []string) float32 {
	hashes := make(chan string, len(posts))
	labels := make(chan bool, len(posts))

	for i := 0; i < 20; i++ {
		go evalWorker(n, hashes, labels)
	}

	for _, p := range posts {
		hashes <- p
	}
	close(hashes)

	trueCounter := 0
	counter := 0
	for i := 0; i < len(posts); i++ {
		if <-labels {
			trueCounter++
		}
		counter++
	}

	return float32(trueCounter) / float32(counter)
}
