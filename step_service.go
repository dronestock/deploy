package main

import (
	"context"
	_ "embed"

	"github.com/goexl/gfx"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/rand"
	"github.com/goexl/gox/tpl"
)

//go:embed template/kubernetes/service.yaml.gohtml
var defaultServiceTemplate []byte

type stepService struct {
	*plugin
}

func newServiceStep(plugin *plugin) *stepService {
	return &stepService{
		plugin: plugin,
	}
}

func (d *stepService) Runnable() bool {
	return nil != d.Service
}

func (d *stepService) Run(_ context.Context) (err error) {
	// 增加端口，兼容只想暴露一个端口的情况
	if 0 != d.Service.Port {
		port := new(servicePort)
		port.Name = d.Name
		port.Port = d.Service.Port
		port.Target = d.Service.Target
		port.Node = d.Service.Node
		port.Protocol = d.Service.Protocol
		d.Service.Ports = append(d.Service.Ports, port)
	}

	label := rand.New().String().Build().Generate()
	filename := gox.StringBuilder(service).Append(dot).Append(label).Append(dot).Append(yaml).String()

	// 写入配置文件
	if "" != d.Kubernetes.Service {
		err = gfx.Copy(d.Service.Template, filename)
	} else {
		err = tpl.New(string(defaultServiceTemplate)).Data(d.plugin).Build().ToFile(filename)
	}
	if nil != err {
		return
	}

	// 清理文件
	d.Cleanup().File(filename).Build()
	err = d.kubectl(args.New().Build().Subcommand(apply).Arg(file, filename).Build())

	return
}
