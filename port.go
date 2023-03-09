package main

import (
	"strings"
)

type port struct {
	// 名称
	Name string `json:"name"`
	// 端口
	Port int `json:"port"`
	// 协议
	Protocol string `default:"tcp" json:"protocol" validate:"oneof=tcp udp sctp"`
}

func (p *port) KubernetesProtocol() string {
	return strings.ToUpper(p.Protocol)
}
