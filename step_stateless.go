package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/goexl/gfx"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/rand"
	"github.com/goexl/gox/tpl"
)

const defaultDeploymentTemplate = "docker/etc/kubernetes/template/deployment.yaml.gohtml"

type stepStateless struct {
	*plugin
}

func newStatelessStep(plugin *plugin) *stepStateless {
	return &stepStateless{
		plugin: plugin,
	}
}

func (s *stepStateless) Runnable() bool {
	return nil != s.Stateless
}

func (s *stepStateless) Run(ctx context.Context) (err error) {
	if nil == s.Stateless.Environments {
		s.Stateless.Environments = make(map[string]string)
	}

	if ee := s.loadKvs(s.filepath(envFilename), s.env); nil != ee {
		err = ee
	} else if ae := s.loadKvs(s.filepath(annotationFilename), s.action); nil != ae {
		err = ae
	}
	if nil != err {
		return
	}

	// 增加端口，兼容只想暴露一个端口的情况
	if 0 != s.Stateless.Port {
		port := new(port)
		port.Name = s.Name
		port.Port = s.Stateless.Port
		port.Protocol = s.Stateless.Protocol
		s.Stateless.Ports = append(s.Stateless.Ports, port)
	}

	// 统一环境变量
	for key, value := range s.Stateless.Envs {
		s.Stateless.Environments[key] = value
	}

	if nil != s.Kubernetes {
		err = s.kubernetes(ctx)
	}

	return
}

func (s *stepStateless) kubernetes(_ context.Context) (err error) {
	if nee := s.loadKvs(s.filepath(s.Kubernetes.Namespace, envFilename), s.env); nil != nee {
		err = nee
	} else if nae := s.loadKvs(s.filepath(s.Kubernetes.Namespace, annotationFilename), s.action); nil != nae {
		err = nae
	}
	if nil != err {
		return
	}

	label := rand.New().String().Build().Generate()
	filename := gox.StringBuilder(stateless).Append(dot).Append(label).Append(dot).Append(yaml).String()
	// 写入配置文件
	if defaultDeploymentTemplate != s.Stateless.Template {
		err = gfx.Copy(s.Stateless.Template, filename)
	} else if bytes, re := os.ReadFile(s.Stateless.Template); nil != re {
		err = re
	} else {
		err = tpl.New(string(bytes)).Data(s.plugin).Build().ToFile(filename)
	}
	if nil != err {
		return
	}

	// 清理文件
	s.Cleanup().File(filename).Build()
	err = s.kubectl(args.New().Build().Subcommand(apply).Arg(file, filename).Build())

	return
}

func (s *stepStateless) filepath(paths ...string) string {
	return filepath.Join(append([]string{s.Dir, stateless}, paths...)...)
}

func (s *stepStateless) env(key string, value string) {
	s.Stateless.Envs[key] = value
}

func (s *stepStateless) action(key string, value string) {
	s.Stateless.Annotations[key] = value
}
