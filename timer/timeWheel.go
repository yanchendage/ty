package timer

import (
	"errors"
	"sync"
	"time"
)

type TimeWheel struct {
	name string
	scales int // 刻度数
	unit int64 // 刻度单位

	index int
	timerQueueSize int
	timerQueue map[int]map[int]*Timer // map[刻度号]map[timer号]timer

	next *TimeWheel //下一层单位较小的时间轮

	sync.RWMutex
}

func (tw *TimeWheel) addTimer(timerID int, timer *Timer) error {

	defer func() error{
		if err := recover(); err != nil{
			return errors.New("addTimer err")
		}
		return  nil
	}()

	interval := timer.callTime - millisecondUnix()

	//间隔大于一个刻度单位
	if interval > tw.unit {
		//几个刻度
		scales := interval / tw.unit
		//比如一圈有12个刻度，当前在刻度3，间隔是10个刻度,那么求余的结果是刻度1
		tw.timerQueue[(tw.index + int(scales)) % tw.scales][timerID] = timer

		return nil
	}

	//间隔小于一个刻度单位，并且没有下一层单位较小的时间轮
	if interval < tw.unit && tw.next == nil{
		//放入当前的刻度
		tw.timerQueue[tw.index][timerID] = timer
		//并执行？
		return nil
	}

	if interval < tw.unit {
		tw.next.addTimer(timerID, timer)
	}

	return nil
}

func (tw *TimeWheel) AddTimer(timerID int, timer *Timer)  {
	//tw.mu.Lock()
	//defer tw.mu.Unlock()
}

func (tw *TimeWheel) run()  {
	for {
		time.Sleep(time.Duration(tw.unit) * time.Millisecond)
		tw.Lock()

		//当前刻度的所有timer
		timers := tw.timerQueue[tw.index]
		tw.timerQueue[tw.index] = make(map[int]*Timer, tw.timerQueueSize)

		for timerID, timer := range timers {

		}


	}
}