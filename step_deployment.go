package main

import (
	"context"
	"os"

	"github.com/goexl/gfx"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/rand"
	"github.com/goexl/gox/tpl"
)

const defaultDeploymentTemplate = "docker/etc/kubernetes/template/deployment.yaml.gohtml"

type stepDeployment struct {
	*plugin
}

func newDeploymentStep(plugin *plugin) *stepDeployment {
	return &stepDeployment{
		plugin: plugin,
	}
}

func (d *stepDeployment) Runnable() bool {
	return nil != d.Deployment
}

func (d *stepDeployment) Run(_ context.Context) (err error) {
	// 增加端口，兼容只想暴露一个端口的情况
	if 0 != d.Deployment.Port {
		port := new(port)
		port.Name = d.Name
		port.Port = d.Deployment.Port
		port.Protocol = d.Deployment.Protocol
		d.Deployment.Ports = append(d.Deployment.Ports, port)
	}

	// 统一环境变量
	if nil == d.Deployment.Environments {
		d.Deployment.Environments = make(map[string]string)
	}
	for key, value := range d.Deployment.Envs {
		d.Deployment.Environments[key] = value
	}

	label := rand.New().String().Build().Generate()
	filename := gox.StringBuilder(deployment).Append(dot).Append(label).Append(dot).Append(yaml).String()

	// 写入配置文件
	if defaultDeploymentTemplate != d.Deployment.Template {
		err = gfx.Copy(d.Deployment.Template, filename)
	} else if bytes, re := os.ReadFile(d.Deployment.Template); nil != re {
		err = re
	} else {
		err = tpl.New(string(bytes)).Data(d.plugin).Build().ToFile(filename)
	}
	if nil != err {
		return
	}

	// 清理文件
	d.Cleanup().File(filename).Build()
	err = d.kubectl(args.New().Build().Subcommand(apply).Arg(file, filename).Build())

	return
}
