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
		for _, hash := range v {
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
			fmt.Println(k, label)
		}
	}
}
