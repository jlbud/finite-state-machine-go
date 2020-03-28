package main

import (
	"fmt"
	"testing"

	"github.com/looplab/fsm"
)

func TestFSM(t *testing.T) {
	// 新构造一个有限状态机
	fsm := fsm.NewFSM(
		"closed", // 初始状态
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},  // 允许从 "closed" 到 "open"
			{Name: "close", Src: []string{"open"}, Dst: "closed"}, // 允许从 "open" 到 "close"
		},
		fsm.Callbacks{},
	)

	// 当前状态
	// 关闭
	fmt.Println(fsm.Current())

	// 从 关闭 到 打开
	err := fsm.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	// 当前状态
	// 打开
	fmt.Println(fsm.Current())

	// 从 打开 到 关闭
	err = fsm.Event("close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm.Current())
}
