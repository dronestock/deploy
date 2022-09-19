package main

type service struct {
	// 名称
	Name string `json:"name" validate:"required"`
	// 端口
	Port int `default:"8080" json:"port" validate:"required"`
}
