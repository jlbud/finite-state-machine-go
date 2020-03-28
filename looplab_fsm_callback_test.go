package main

import (
	"fmt"
	"testing"

	"github.com/looplab/fsm"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
		},
	)

	return d
}

// door实现的动作，会通过相关状态触发
func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

func TestFSMCallback(t *testing.T) {
	door := NewDoor("heaven")

	// 触发动作，改变状态的同时，
	// 会触发回调事件
	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	// 触发动作，改变状态的同时，
	// 会触发回调事件
	err = door.FSM.Event("close")
	if err != nil {
		fmt.Println(err)
	}
}
