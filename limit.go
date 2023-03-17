package main

import (
	"github.com/goexl/gox"
)

type limit struct {
	// 核数
	Cpu float32 `json:"cpu"`
	// 内存
	Memory gox.Size `json:"memory"`
}
