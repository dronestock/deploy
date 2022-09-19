package main

import (
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type plugin struct {
	drone.Base

	// 类型
	Type typ `default:"${PLUGIN_TYPE=${TYPE=deployment}}" validate:"required,oneof=deployment,stateful"`
	// 名称
	Name string `default:"${PLUGIN_NAME=${NAME}}" validate:"required"`
	// 镜像名称
	Image string `default:"${PLUGIN_IMAGE=${IMAGE}}" validate:"required"`

	// 服务配置
	Service *service `default:"${PLUGIN_SERVICE=${SERVICE}}"`
}

func newPlugin() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Steps() []*drone.Step {
	return []*drone.Step{
		drone.NewStep(p.build, drone.Name(`打包`)),
	}
}

func (p *plugin) Fields() (fields gox.Fields) {
	fields = gox.Fields{}
	if nil != p.Service {
		fields = append(fields, field.String(`service.name`, p.Service.Name))
		fields = append(fields, field.Int(`service.port`, p.Service.Port))
	}

	return
}
