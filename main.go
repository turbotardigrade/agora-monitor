package main

import (
	"fmt"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"os"
	"sort"
	"time"
)

const MyNodePath = "./data/monitorNode"

var NodeList = []string{
	"QmWSEKJqQqBNCCHEy4WTZcuZUmjmS2cexnRofceZf2bD7D", //spam
	"QmajBiibByrMkaGFznNrA5x9oxziBMJ5QXsH5Pb8bcPpsS", // spam
	"Qmae3pLmGUJhMMP8ungzHLnA5BYJWzFhbkMxHveC4CVcUU",
	"QmTDhg18AZV27haxPW7ZZPzbczpiRBr2u7ZVznaA3Eu13J",
	"Qmac7EgAPVuGcJcuBvD51Jpt8TkM7pMfQz1TjXi1pfgCzC",
	"QmYSmjBiEQdkU56dkrb1cvNLfPAGKoRFXXxUAqScbEZUT8",
	"QmaUFw1RWSzgeV1qNBDKx8okHTBEUqUiWrxKSZZjXBMYtB",
	"QmSvFuDHQCcGiQrNXaDeMn8BSTeDgDYs5Pg3vsbbMRgyS3",
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

	for {
		monitorRoutine(n, f)
		time.Sleep(5 * time.Second)
	}
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
