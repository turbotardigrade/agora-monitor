package node

import (
	"bytes"
	peer "gx/ipfs/QmWUswjn261LSyVxWAEpMVtPdy8zmKBJJfBpG3Qdpa8ZsE/go-libp2p-peer"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/corenet"
)

func Request(node *core.IpfsNode, targetPeer string, path string, body string) (string, error) {
	// Check if Node hash is valid
	target, err := peer.IDB58Decode(targetPeer)
	if err != nil {
		return "", err
	}

	// Connect to targetPeer
	stream, err := corenet.Dial(node, target, path)
	if err != nil {
		return "", err
	}

	// Exchange request and response
	stream.Write([]byte(body))

	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String(), nil
}
