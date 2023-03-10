package main

type servicePort struct {
	port

	// 目标端口
	Target int `json:"target"`
	// 节点端口
	Node int `json:"node"`
}
