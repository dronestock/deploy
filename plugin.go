package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
)

type plugin struct {
	drone.Base

	// Kubernetes配置
	Kubernetes *kubernetes `default:"${KUBERNETES}" json:"kubernetes"`

	// 用户名
	Username string `default:"${USERNAME=default}" json:"username"`
	// 密钥
	Token string `default:"${TOKEN}" json:"token" validate:"required_without_all=Key Password"`
	// 密码
	// 密钥和密码统一使用密码做内部变量配置
	Password string `default:"${PASSWORD}" json:"password" validate:"required_without_all=Key Token"`
	// 密钥
	Key string `default:"${KEY}" json:"key" validate:"required_without_all=Token Password"`

	// 名称
	Name string `default:"${NAME}" json:"name" validate:"required"`
	// 注册表
	Registry string `default:"${REGISTRY}" json:"registry" validate:"required"`
	// 仓库
	Repository string `default:"${REPOSITORY}" json:"repository" validate:"required"`
	// 标签
	Tag string `default:"${TAG=latest}" json:"tag"`

	// 配置目录
	Dir string `default:"${DIR=.deploy}" json:"dir"`
	// 无状态服务
	Stateless *_stateless `default:"${DEPLOYMENT}" json:"stateless"`
	// 服务配置
	Service *_service `default:"${SERVICE}" json:"service"`
}

func newPlugin() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Setup() (err error) {
	p.Password = gox.If("" != p.Token, p.Token)

	return
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(newSetupStep(p)).Name("配置").Build(),
		drone.NewStep(newNamespaceStep(p)).Name("命名空间").Build(),
		drone.NewStep(newStatelessStep(p)).Name("无状态应用").Build(),
		drone.NewStep(newServiceStep(p)).Name("服务").Build(),
	}
}

func (p *plugin) Fields() (fields gox.Fields[any]) {
	fields = make(gox.Fields[any], 0, 2)
	if nil != p.Stateless {
		fields = append(fields, field.New("stateless", p.Stateless))
	}
	if nil != p.Service {
		fields = append(fields, field.New("service", p.Service))
	}

	return
}

func (p *plugin) kubectl(args *args.Args) (err error) {
	_, err = p.Command(p.Kubernetes.Binary).Args(args).Build().Exec()

	return
}

func (p *plugin) outputs(args *args.Args, outputs *[]string) (err error) {
	_, err = p.Command(p.Kubernetes.Binary).Args(args).Collector().TrimRight(enter).Strings(outputs).Build().Exec()

	return
}

func (p *plugin) loadKvs(path string, fun kvFun) (err error) {
	if _, se := os.Stat(path); nil == se {
		err = p.readline(path, p.kv(fun))
	}

	return
}

func (p *plugin) kv(fun kvFun) lineFun {
	return func(line string) {
		values := strings.Split(line, equal)
		if 2 == len(values) {
			fun(values[0], values[1])
		}
	}
}

func (p *plugin) readline(path string, fun lineFun) (err error) {
	if file, oe := os.Open(path); nil != oe {
		err = oe
	} else {
		defer p.close(file)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fun(scanner.Text())
		}
		err = scanner.Err()
	}

	return
}

func (p *plugin) close(file *os.File) {
	_ = file.Close()
}
