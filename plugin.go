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
	Kubernetes *kubernetes `default:"${KUBERNETES}" json:"kubernetes,omitempty"`

	// 用户名
	Username string `default:"${USERNAME=default}" json:"username,omitempty"`
	// 密钥
	Token string `default:"${TOKEN}" json:"token,omitempty" validate:"required_without_all=Key Password"`
	// 密码
	// 密钥和密码统一使用密码做内部变量配置
	Password string `default:"${PASSWORD}" json:"password,omitempty" validate:"required_without_all=Key Token"`
	// 密钥
	Key string `default:"${KEY}" json:"key,omitempty" validate:"required_without_all=Token Password"`

	// 名称
	Name string `default:"${NAME}" json:"name,omitempty" validate:"required"`
	// 注册表
	Registry string `default:"${REGISTRY}" json:"registry,omitempty" validate:"required"`
	// 仓库
	Repository string `default:"${REPOSITORY}" json:"repository,omitempty" validate:"required"`
	// 标签
	Tag string `default:"${TAG=latest}" json:"tag,omitempty"`

	// 端口
	Port *port `default:"${PORT}" json:"port,omitempty"`
	// 端口列表
	Ports []*port `default:"${PORTS}" json:"ports,omitempty"`

	// 资源限制
	Resource *resource `default:"${resource}" json:"resource,omitempty"`
	// 配置目录
	Dir string `default:"${DIR=.deploy}" json:"dir,omitempty"`
	// 无状态服务
	Stateless *_stateless `default:"${STATELESS}" json:"stateless,omitempty"`

	// 注解
	Annotations map[string]string `default:"${ANNOTATIONS}" json:"annotations,omitempty"`
	// 环境变量
	Environments map[string]string `default:"${ENVIRONMENTS}" json:"environments,omitempty"`
}

func newPlugin() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Setup() (err error) {
	p.Password = gox.If("" != p.Token, p.Token)
	// 统一端口配置
	if nil == p.Ports {
		p.Ports = make([]*port, 0)
	}
	if nil != p.Port {
		p.Ports = append(p.Ports, p.Port)
	}

	// 统一注解
	if 0 != len(p.Annotations) && nil != p.Stateless {
		p.Stateless.Annotations = gox.If(nil == p.Stateless.Annotations, make(map[string]string))
		p.copy(p.Annotations, p.Stateless.Annotations)
	}
	// 统一环境变量
	if 0 != len(p.Environments) && nil != p.Stateless {
		p.Stateless.Environments = gox.If(nil == p.Stateless.Environments, make(map[string]string))
		p.copy(p.Environments, p.Stateless.Environments)
	}

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
	fields = gox.Fields[any]{
		field.New("username", p.Username),
		field.New("name", p.Name),
		field.New("registry", p.Registry),
		field.New("repository", p.Repository),
		field.New("tag", p.Tag),
		field.New("ports", p.Ports),
		field.New("dir", p.Dir),
	}
	if nil != p.Stateless {
		fields.Add(field.New("stateless", p.Stateless))
	}
	if nil != p.Resource {
		fields.Add(field.New("resource", p.Resource))
	}
	if 0 != len(p.Annotations) {
		fields.Add(field.New("annotations", p.Annotations))
	}
	if 0 != len(p.Environments) {
		fields.Add(field.New("environments", p.Annotations))
	}

	return
}

func (p *plugin) KubernetesRepository() string {
	return strings.ReplaceAll(p.Repository, slash, dot)
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

func (p *plugin) copy(from map[string]string, to map[string]string) {
	for key, value := range from {
		to[key] = value
	}
}
