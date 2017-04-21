package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/repo/config"
)

const MyNodePath = "./data/monitorNode"

var NodeList = []Node{
	{"QmeMVVuVKEWs4SM8tpc2DPt8vbpxdQwRokBwXSu5putT4Z", false},
	{"QmUXd9sH7dvNcJwWQwAdX1CKZoFGFrrvHSoPisfjH7fubE", false},
	{"QmWtXyS6AcEBFyLCNcEbr8rKdga5toFSAQFNpDM5ZTkvpq", false},
	{"QmbWSkBe2hLcksxQ7dyWNZ9q33fYUwzeeQyiKk6mXdiRew", false},
	{"QmT7LNeXzHKWRSAZq3U5NHcTsPUdR9my24gu29LVX7MgsS", false},
	{"QmQ9epk6VuJK7aJVLdKkGrdjTNDZXu2CC71uKcYwRAvFni", true},
	{"Qmcpr5tJqM8uVP2TV9Di6KQVQ6Dbo9hq9Y83osTzwKvj8f", true},
	{"QmZVRixZR14rYvkmsjpmSTBs8TpGy6wu1joUHT7eAfZupU", true},
}

func main() {
	n, err := InitNode()
	if err != nil {
		panic(err)
	}

	err = loadLabels()
	if err != nil {
		panic(err)
	}

	err = CreateFileIfNotExists("stats.csv")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("stats.csv", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	go func() {
		for {
			monitorRoutine(n, f)
		}
	}()

	http.HandleFunc("/nodes", nodesHandler)
	http.HandleFunc("/monitor", monitorHandler)
	http.Handle("/", http.FileServer(http.Dir("visual")))

	fmt.Println("Monitor API running")
	http.ListenAndServe(":8080", nil)
}

func monitorRoutine(n *core.IpfsNode, f *os.File) {
	healthy, posts, blacklists := monitor(n)
	sortedList := sortedNodes(healthy)

	fmt.Println("----------------------------------------------------------------------")
	for _, h := range sortedList {
		if !healthy[h] {
			fmt.Println(formatHash(h), "is unhealthy")
		}
	}

	fmt.Println("\nPeer\t\tnum_posts\tspamratio\tblacklist_count")
	hr, min, sec := time.Now().Clock()
	line := fmt.Sprintf("%d:%d:%d,", hr, min, sec)
	for _, h := range sortedList {

		time.Sleep(1 * time.Second)

		ps := posts[h]
		total := len(ps)
		spamratio := evaluatePosts(n, ps)
		blacklistCount := len(blacklists[h])

		MonitorMapLock.Lock()
		if healthy[h] {
			MonitorMap[h] = MonitorResp{
				healthy[h],
				spamratio,
				blacklists[h],
			}
		} else {
			resp := MonitorResp{}
			resp.Healthy = healthy[h]
			MonitorMap[h] = resp
		}
		MonitorMapLock.Unlock()

		fmt.Printf("%v  \t%0000d    \t\t%3.2f   \t\t%v\n", formatHash(h), total, spamratio, blacklistCount)
		line += fmt.Sprintf("%v,%v,%v,%v,", h, total, spamratio, blacklistCount)
	}
	line += "\n"

	fmt.Println("\nBlacklists")
	for _, h := range sortedList {
		fmt.Print(formatHash(h) + ":  ")
		for _, b := range blacklists[h] {
			for k, _ := range b {
				fmt.Print(formatHash(k))
			}
		}
		fmt.Println("")
	}

	_, err := f.WriteString(line)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func sortedNodes(nodes map[string]bool) []string {
	arr := make([]string, len(nodes))

	i := 0
	for k, _ := range nodes {
		arr[i] = k
		i += 1
	}

	sort.Strings(arr)

	return arr
}
func formatHash(hash string) string {
	return "[" + hash[len(hash)-5:len(hash)] + "]\t"
}

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
		err := NewNodeRepo(MyNodePath, addr)
		if err != nil {
			return nil, err
		}

		fmt.Println("Seeding...")
		time.Sleep(5 * time.Second)
		fmt.Println("Seeding done.")

	}

	return NewNode(MyNodePath)
}
