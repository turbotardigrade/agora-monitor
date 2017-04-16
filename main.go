package main

import (
	"fmt"
	"sort"
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

	healthy, posts := monitor(n)
	sortedList := sortedNodes(healthy)

	fmt.Println("Health Status")
	for _, h := range sortedList {
		fmt.Println(h, healthy[h])
	}

	fmt.Println("\nPosts")
	for _, h := range sortedList {
		ps := posts[h]
		total := len(ps)

		fmt.Println("\nEvaluate peer ", h)
		fmt.Println("Total:", total)
		fmt.Println("Spam ratio:", evaluatePosts(n, ps))
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
