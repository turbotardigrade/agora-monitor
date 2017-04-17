package main

import (
	"fmt"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"net/http"
	"os"
	"sort"
	"time"
)

const MyNodePath = "./data/monitorNode"

var NodeList = []Node{
	{"QmRHbpV4aEVVhoJ5F9dvRpxULLc1nDdQf6ZQqHv1pWRGrs", true}, // spammer
	// {"QmajBiibByrMkaGFznNrA5x9oxziBMJ5QXsH5Pb8bcPpsS", true}, // spammer
	// {"Qmae3pLmGUJhMMP8ungzHLnA5BYJWzFhbkMxHveC4CVcUU", false},
	// {"QmTDhg18AZV27haxPW7ZZPzbczpiRBr2u7ZVznaA3Eu13J", false},
	// {"Qmac7EgAPVuGcJcuBvD51Jpt8TkM7pMfQz1TjXi1pfgCzC", false},
	// {"QmYSmjBiEQdkU56dkrb1cvNLfPAGKoRFXXxUAqScbEZUT8", false},
	// {"QmaUFw1RWSzgeV1qNBDKx8okHTBEUqUiWrxKSZZjXBMYtB", false},
	{"QmcpY55BvoQWQxSqotpTWAr5DG14mhsgcBFchRXjYoSuBu", false},
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
			time.Sleep(5 * time.Second)
		}
	}()

	http.HandleFunc("/nodes", nodesHandler)
	http.HandleFunc("/monitor", monitorHandler)
	http.Handle("/", http.FileServer(http.Dir("visual")))

	http.ListenAndServe(":8080", nil)
	fmt.Println("Monitor API running")
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

	fmt.Println("\nPosts")
	hr, min, sec := time.Now().Clock()
	line := fmt.Sprintf("%d:%d:%d,", hr, min, sec)
	for _, h := range sortedList {
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

		fmt.Printf("%v  %v\t%v\t%v\n", formatHash(h), total, spamratio, blacklistCount)
		line += fmt.Sprintf("%v,%v,%v,%v,", h, total, spamratio, blacklistCount)
	}
	line += "\n"

	fmt.Println("\nBlacklists")
	for _, h := range sortedList {
		fmt.Println(formatHash(h), blacklists[h])
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
