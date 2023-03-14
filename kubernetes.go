package main

type kubernetes struct {
	// 控制程序
	Binary string `default:"${KUBERNETES_BINARY=kubectl}" json:"binary"`
	// 服务
	Endpoint string `default:"${KUBERNETES_ENDPOINT}" validate:"required" json:"endpoint"`
	// 命名空间
	Namespace string `default:"${KUBERNETES_NAMESPACE=default}" json:"namespace"`
	// 共享
	Share share `default:"${KUBERNETES_SHARE}" json:"share"`
	// 版本
	Version string `default:"${KUBERNETES_VERSION=v1}"`
}
