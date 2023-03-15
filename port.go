package main

import (
	"strings"
)

type port struct {
	// 名称
	Name string `json:"name,omitempty" validate:"required"`
	// 本地端口
	Local int `json:"local,omitempty" validate:"required"`
	// 暴露端口
	Expose int `json:"expose,omitempty"`
	// 协议
	Protocol string `default:"tcp" json:"protocol,omitempty" validate:"oneof=tcp udp sctp"`
}

func (p *port) KubernetesProtocol() string {
	return strings.ToUpper(p.Protocol)
}
