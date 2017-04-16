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
	"Qmdwu4tKzVLhkytBM4gLuDUNR17sYdEMPY7hTfANe5H2r8", //Left spam
	"Qmc3R6A1Dd9RtQdMqw5udPihjy1gsKs38imiPqfHZ49ZPb", //Right Top
	"QmcDFEbAQtqL4Dh2xAZ4GTqPHUuAdhnvLGAJfgpH1Q1B2t", //Right Bot
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
	healthy, posts := monitor(n)
	sortedList := sortedNodes(healthy)

	fmt.Println("----------------------------------------------------------------------")
	for _, h := range sortedList {
		if !healthy[h] {
			fmt.Println(formatHash(h), "is unhealthy")
		}
	}

	fmt.Println("\nPosts")
	line := ""
	for _, h := range sortedList {
		ps := posts[h]
		total := len(ps)
		spamratio := evaluatePosts(n, ps)

		fmt.Println(formatHash(h), total, spamratio)
		line += fmt.Sprintf("%v,%v,%v,", h, total, spamratio)
	}

	line += "\n"

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
