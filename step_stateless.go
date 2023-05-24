package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/goexl/gfx"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/rand"
	"github.com/goexl/gox/tpl"
)

//go:embed template/kubernetes/deployment.yaml.gohtml
var defaultDeploymentTemplate []byte

type stepStateless struct {
	*plugin

	printed bool
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
	// 保证在插件退时清理文件
	s.Cleanup().File(filename).Build()

	// 写入配置文件
	if "" != s.Kubernetes.Deployment {
		err = gfx.Copy(s.Kubernetes.Deployment, filename)
	} else {
		err = tpl.New(string(defaultDeploymentTemplate)).Data(s.plugin).Build().ToFile(filename)
	}
	if nil != err {
		return
	}

	ka := args.New().Build()
	ka.Subcommand(apply)
	ka.Arg(file, filename)
	ka.Arg(force, strconv.FormatBool(true))
	if err = s.kubectl(ka.Build()); nil != err && !s.printed {
		bytes, _ := os.ReadFile(filename)
		fmt.Println(string(bytes))
		s.printed = true
	}

	return
}

func (s *stepStateless) filepath(paths ...string) string {
	return filepath.Join(append([]string{s.Dir, stateless}, paths...)...)
}

func (s *stepStateless) env(key string, value string) {
	s.Stateless.Environments[key] = gox.StringBuilder(quota, value, quota).String()
}

func (s *stepStateless) action(key string, value string) {
	s.Stateless.Annotations[key] = value
}
