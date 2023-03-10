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

const defaultServiceTemplate = "docker/etc/kubernetes/template/_service.yaml.gohtml"

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
	if defaultServiceTemplate != d.Service.Template {
		err = gfx.Copy(d.Service.Template, filename)
	} else if bytes, re := os.ReadFile(d.Service.Template); nil != re {
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
