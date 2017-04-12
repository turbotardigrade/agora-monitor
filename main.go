package main

import (
	"fmt"
)

const MyNodePath = "./data/monitorNode"

var NodeList = []string{
	"QmRJe2QeEt89qeEM4onq7AmgKkLNEwJuWsRoR97Zdnx17C", //Left spam
	"QmPrEnycMnzg6sADkcPRLefYWf7bnb8TG6k13nsivV6noX", //Right Top
	"QmWCBPbwi9JCRAG1AE9ik2E3Y3FD3Xc2ARj6k6QREZdGBy", //Right Bot
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

	fmt.Println("Health Status")
	for k, v := range healthy {
		fmt.Println(k, v)
	}

	fmt.Println("\nPosts")
	for k, v := range posts {
		fmt.Println("\nEvaluate peer ", k)
		fmt.Println("Total:", len(v))
		fmt.Println("Spam ratio:", evaluatePosts(n, v))
	}
}
