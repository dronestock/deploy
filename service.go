package main

type _service struct {
	// 模板
	// nolint: lll
	Template string `default:"${DEPLOYMENT_TEMPLATE=docker/etc/kubernetes/template/_service.yaml.gohtml}" json:"template,omitempty"`

	// 端口
	Port int `default:"8080" json:"port"`
	// 协议
	Protocol string `default:"tcp" json:"protocol" validate:"oneof=tcp udp sctp"`
	// 目标端口
	Target int `json:"target"`
	// 节点端口
	Node int `json:"node"`

	// 端口列表
	Ports []*servicePort `json:"ports"`
}
