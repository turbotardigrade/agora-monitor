package main

import (
	"fmt"
)

const MyNodePath = "./data/monitorNode"

var NodeList = []string{
	"QmdtfJBMitotUWBX5YZ6rYeaYRFu6zfXXMZP6fygEWK2iu",
	"QmVmPkKN9XXfxwQfinSWDYuU8M6U9dZdL46uSoSwuYgLKL",
}

func main() {
	n, err := InitNode()
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
		fmt.Println(k, v)
	}
}
