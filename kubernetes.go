package main

type kubernetes struct {
	// 控制程序
	Binary string `default:"${KUBERNETES_BINARY=kubectl}" json:"binary"`
	// 无状态服务模板
	Deployment string `default:"${KUBERNETES_DEPLOYMENT}" json:"deployment,omitempty"`
	// 服务模板
	Service string `default:"${KUBERNETES_SERVICE}" json:"service,omitempty"`
	// 服务
	Endpoint string `default:"${KUBERNETES_ENDPOINT}" validate:"required" json:"endpoint"`
	// 命名空间
	Namespace string `default:"${KUBERNETES_NAMESPACE=default}" json:"namespace"`
	// 共享
	Share share `default:"${KUBERNETES_SHARE}" json:"share"`
}
