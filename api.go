package main

import (
	"encoding/json"
	"net/http"
)

type NodesResp struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	ID        string `json:"id"`
	IsSpammer bool   `json:"is_spammer"`
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes := []Node{}
	for _, n := range NodeList {
		nodes = append(nodes, n)
	}

	js, err := json.Marshal(nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type MonitorResp struct {
	Healthy   bool                 `json:"healthy"`
	SpamRatio float32              `json:"spam_ratio"`
	Blacklist []map[string]float32 `json:"blacklist"`
}

func monitorHandler(w http.ResponseWriter, r *http.Request) {
	MonitorMapLock.RLock()
	js, err := json.Marshal(MonitorMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	MonitorMapLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
