package main

type _deployment struct {
	// 模板
	Template string `default:"${DEPLOYMENT_TEMPLATE=docker/etc/kubernetes/template/deployment.yaml.gohtml}"`

	// 端口
	Port int `default:"8080" json:"port"`
	// 协议
	Protocol string `default:"tcp" json:"protocol" validate:"oneof=tcp udp"`

	// 端口列表
	Ports []*port `json:"ports"`
	// 复本数
	Replicas int `json:"replicas" validate:"required"`
	// 注解
	Annotations map[string]string `json:"annotations"`
	// 环境变量
	Envs map[string]string `json:"envs"`
	// 环境变量
	Environments map[string]string `json:"environment"`
}
