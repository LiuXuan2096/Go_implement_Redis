package database

import "strings"

// 存放 指令 和指令对应的执行函数的容器
// key是指令 value是 command
var cmdTable = make(map[string]*command)

// command 有两个成员 1. 指令对应的执行函数 2.指令对应的参数数量用于参数校验
type command struct {
	executor ExecFunc
	arity    int
}

// RegisterCommand 向 cmdTable 中注册指令和该指令对应的 command 结构体变量
func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
