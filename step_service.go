package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

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

	printed bool
}

func newServiceStep(plugin *plugin) *stepService {
	return &stepService{
		plugin: plugin,
	}
}

func (s *stepService) Runnable() (runnable bool) {
	for _, port := range s.Ports {
		if 0 != port.Expose {
			runnable = true
		}

		if runnable {
			break
		}
	}

	return
}

func (s *stepService) Run(ctx context.Context) (err error) {
	if nil != s.Kubernetes {
		err = s.kubernetes(ctx)
	}

	return
}

func (s *stepService) kubernetes(_ context.Context) (err error) {
	label := rand.New().String().Build().Generate()
	filename := gox.StringBuilder(service).Append(dot).Append(label).Append(dot).Append(yaml).String()

	// 写入配置文件
	if "" != s.Kubernetes.Service {
		err = gfx.Copy(s.Kubernetes.Service, filename)
	} else {
		err = tpl.New(string(defaultServiceTemplate)).Data(s.plugin).Build().ToFile(filename)
	}
	if nil != err {
		return
	}

	// 清理文件
	s.Cleanup().File(filename).Build()
	if err = s.kubectl(args.New().Build().Subcommand(apply).Arg(file, filename).Build()); nil != err && !s.printed {
		bytes, _ := os.ReadFile(filename)
		fmt.Println(string(bytes))
		s.printed = true
	}

	return
}
