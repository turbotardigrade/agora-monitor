package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core"
	"gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/core/coreunix"
)

var labels = make(map[string]bool)

func loadLabels() error {
	file, err := ioutil.ReadFile("./labels.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &labels)
	if err != nil {
		return err
	}

	fmt.Println(len(labels))
	return nil
}

func checkLabel(content string) (bool, error) {
	h := md5.New()
	io.WriteString(h, content)

	b, ok := labels[hex.EncodeToString(h.Sum(nil))]
	if !ok {
		return false, errors.New("Content is not labeled")
	}
	return b, nil
}

func getContent(n *core.IpfsNode, hash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := coreunix.Cat(ctx, n, hash)
	if err != nil {
		return "", err
	}

	obj := make(map[string]interface{})
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(buf, &obj)
	if err != nil {
		return "", err
	}

	data, ok := obj["Data"].(map[string]interface{})
	if !ok {
		return "", errors.New("IPFS obj format broken")
	}

	content, ok := data["Content"].(string)
	if !ok {
		return "", errors.New("IPFS obj format broken")
	}

	return content, nil
}
