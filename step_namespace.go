package main

import (
	"context"

	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
)

const (
	namespace = "namespace"
	noHeaders = "no-headers"
	output    = "output"
	onlyName  = "custom-columns=:metadata.name"
)

type stepNamespace struct {
	*plugin
}

func newNamespaceStep(plugin *plugin) *stepNamespace {
	return &stepNamespace{
		plugin: plugin,
	}
}

func (n *stepNamespace) Runnable() bool {
	return nil != n.Kubernetes
}

func (n *stepNamespace) Run(_ context.Context) (err error) {
	// 查看现在的所有命名空间
	getArgs := args.New().Build()
	getArgs.Add(get, namespace).Flag(noHeaders).Flag(output).Add(onlyName)
	namespaces := make([]string, 0, 1)
	if err = n.outputs(getArgs.Build(), &namespaces); nil != err {
		return
	}

	// 如果命名空间不存在，创建命名空间
	if !gox.Contains(&namespaces, n.Kubernetes.Namespace) {
		err = n.kubectl(args.New().Build().Subcommand(create, namespace, n.Kubernetes.Namespace).Build())
	}

	return
}
