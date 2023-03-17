package main

type resource struct {
	// 限制
	Limit *limit `json:"limit"`
	// 请求
	Request *limit `json:"request"`
}
