package main

type _stateless struct {
	// 模板
	// nolint: lll
	Template string `default:"${DEPLOYMENT_TEMPLATE=docker/etc/kubernetes/template/deployment.yaml.gohtml}" json:"template,omitempty"`

	// 端口
	Port int `default:"8080" json:"port,omitempty"`
	// 协议
	Protocol string `default:"tcp" json:"protocol,omitempty" validate:"oneof=tcp udp sctp"`

	// 端口列表
	Ports []*port `json:"ports,omitempty"`
	// 复本数
	Replicas int `json:"replicas,omitempty" validate:"required"`
	// 注解
	Annotations map[string]string `json:"annotations,omitempty"`
	// 环境变量
	Envs map[string]string `json:"envs,omitempty"`
	// 环境变量
	Environments map[string]string `json:"environment,omitempty"`
}
