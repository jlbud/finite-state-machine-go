// main.go
package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//////////////////// 状态实现开始
type FSMState string            // 状态
type FSMEvent string            // 事件
type FSMHandler func() FSMState // 处理方法，并返回新的状态

// 有限状态机
type FSM struct {
	mu       sync.Mutex                           // 排他锁
	state    FSMState                             // 当前状态
	handlers map[FSMState]map[FSMEvent]FSMHandler // 处理地图集，每一个状态都可以出发有限个事件，执行有限个处理
}

// 获取当前状态
func (f *FSM) getState() FSMState {
	return f.state
}

// 设置当前状态
func (f *FSM) setState(newState FSMState) {
	f.state = newState
}

// 某状态添加事件处理方法
func (f *FSM) AddHandler(state FSMState, event FSMEvent, handler FSMHandler) *FSM {
	if _, ok := f.handlers[state]; !ok {
		f.handlers[state] = make(map[FSMEvent]FSMHandler)
	}
	if _, ok := f.handlers[state][event]; ok {
		fmt.Printf("[警告] 状态(%s)事件(%s)已定义过", state, event)
	}
	f.handlers[state][event] = handler
	return f
}

// 事件处理
func (f *FSM) Call(event FSMEvent) FSMState {
	f.mu.Lock()
	defer f.mu.Unlock()
	events := f.handlers[f.getState()]
	if events == nil {
		return f.getState()
	}
	if fn, ok := events[event]; ok {
		oldState := f.getState()
		f.setState(fn())
		newState := f.getState()
		fmt.Println("状态从 [", oldState, "] 变成 [", newState, "]")
		fmt.Println()
	}
	return f.getState()
}

// 实例化FSM
func NewFSM(initState FSMState) *FSM {
	return &FSM{
		state:    initState,
		handlers: make(map[FSMState]map[FSMEvent]FSMHandler),
	}
}

//////////////////// 状态实现结束

var (
	Poweroff        = FSMState("关闭")
	FirstGear       = FSMState("1档")
	SecondGear      = FSMState("2档")
	ThirdGear       = FSMState("3档")
	PowerOffEvent   = FSMEvent("按下关闭按钮")
	FirstGearEvent  = FSMEvent("按下1档按钮")
	SecondGearEvent = FSMEvent("按下2档按钮")
	ThirdGearEvent  = FSMEvent("按下3档按钮")
	PowerOffHandler = FSMHandler(func() FSMState {
		fmt.Println("电风扇已关闭")
		return Poweroff
	})
	FirstGearHandler = FSMHandler(func() FSMState {
		fmt.Println("电风扇开启1档，微风徐来！")
		return FirstGear
	})
	SecondGearHandler = FSMHandler(func() FSMState {
		fmt.Println("电风扇开启2档，凉飕飕！")
		return SecondGear
	})
	ThirdGearHandler = FSMHandler(func() FSMState {
		fmt.Println("电风扇开启3档，发型被吹乱了！")
		return ThirdGear
	})
)

// 电风扇
type ElectricFan struct {
	*FSM
}

// 实例化电风扇
func NewElectricFan(initState FSMState) *ElectricFan {
	return &ElectricFan{
		FSM: NewFSM(initState),
	}
}

// 入口函数
func TestSelfFSM(t *testing.T) {

	efan := NewElectricFan(Poweroff) // 初始状态是关闭的
	// 关闭状态
	efan.AddHandler(Poweroff, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(Poweroff, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(Poweroff, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(Poweroff, ThirdGearEvent, ThirdGearHandler)
	// 1档状态
	efan.AddHandler(FirstGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(FirstGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(FirstGear, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(FirstGear, ThirdGearEvent, ThirdGearHandler)
	// 2档状态
	efan.AddHandler(SecondGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(SecondGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(SecondGear, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(SecondGear, ThirdGearEvent, ThirdGearHandler)
	// 3档状态
	efan.AddHandler(ThirdGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(ThirdGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(ThirdGear, SecondGearEvent, SecondGearHandler)
	efan.AddHandler(ThirdGear, ThirdGearEvent, ThirdGearHandler)

	// 开始测试状态变化
	efan.Call(ThirdGearEvent) // 按下3档按钮
	efan.Call(PowerOffEvent)  // 按下关闭按钮
	efan.Call(FirstGearEvent) // 按下1档按钮
	// golang 1.14新特性，会边执行边输出，
	// 1.14前，只会最终执行完整体输出
	time.Sleep(1 * time.Second)
	efan.Call(SecondGearEvent) // 按下2档按钮
	efan.Call(PowerOffEvent)   // 按下关闭按钮
}
