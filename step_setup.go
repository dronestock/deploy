package main

import (
	"context"

	"github.com/goexl/gox/arg"
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
	return true
}

func (s *stepSetup) Run(_ context.Context) (err error) {
	// 设置密钥
	tokenArgs := arg.New().Build()
	tokenArgs.Add(config, setCredentials, _default).Long(token, s.Token)
	if err = s.kubectl(tokenArgs.Build()); nil != err {
		return
	}

	// 设置通信服务器
	serverArgs := arg.New().Build()
	serverArgs.Add(config, setCluster, _default).Long(server, s.Endpoint).Long(insecure, true)
	if err = s.kubectl(serverArgs.Build()); nil != err {
		return
	}

	// 设置用户
	userArgs := arg.New().Build()
	userArgs.Add(config, setContext, _default).Long(cluster, _default).Long(user, s.Username)
	if err = s.kubectl(userArgs.Build()); nil != err {
		return
	}

	// 设置上下文
	contextArgs := arg.New().Build()
	contextArgs.Add(config, setContext, useContext, _default)
	err = s.kubectl(contextArgs.Build())

	return
}
