package main

import (
	"context"

	"github.com/goexl/gox/args"
)

const (
	setCredentials = "set-credentials"
	setContext     = "set-context"
	setCluster     = "set-cluster"
	server         = "server"
	cluster        = "cluster"
	insecure       = "insecure-skip-tls-verify"
	user           = "user"
	useContext     = "use-context"
)

type stepSetup struct {
	*plugin
}

func newSetupStep(plugin *plugin) *stepSetup {
	return &stepSetup{
		plugin: plugin,
	}
}

func (s *stepSetup) Runnable() bool {
	return nil != s.Kubernetes
}

func (s *stepSetup) Run(ctx context.Context) (err error) {
	if nil != s.Kubernetes {
		err = s.kubernetes(ctx)
	}

	return
}

func (s *stepSetup) kubernetes(_ context.Context) (err error) {
	// 设置密钥
	tokenArgs := args.New().Build()
	tokenArgs.Subcommand(config, setCredentials, def).Arg(token, s.Password)
	if err = s.kubectl(tokenArgs.Build()); nil != err {
		return
	}

	// 设置通信服务器
	serverArgs := args.New().Build()
	serverArgs.Subcommand(config, setCluster, def).Arg(server, s.Kubernetes.Endpoint).Arg(insecure, true)
	if err = s.kubectl(serverArgs.Build()); nil != err {
		return
	}

	// 设置用户
	userArgs := args.New().Build()
	userArgs.Subcommand(config, setContext, def).Arg(cluster, def).Arg(user, s.Username)
	if err = s.kubectl(userArgs.Build()); nil != err {
		return
	}

	// 设置上下文
	contextArgs := args.New().Build()
	contextArgs.Subcommand(config, useContext, def)
	err = s.kubectl(contextArgs.Build())

	return
}
