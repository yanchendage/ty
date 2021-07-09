package timer

import "time"

type Timer struct {
	task *Task //任务
	callTime int64 //调用时间
}

//返回1970-1-1至今的毫秒数
func millisecondUnix() int64 {
	return time.Now().UnixNano()  / 1e6
}

//绝对时间
func NewAbsoluteTimer(task *Task, absoluteTime int64)  *Timer {
	return &Timer{
		task:    task,
		callTime: absoluteTime / 1e6,
	}
}

//相对时间
func NewRelativeTimer(task *Task, relativeTime time.Duration)  *Timer  {
	return NewAbsoluteTimer(task, time.Now().UnixNano()+int64(relativeTime))
}

//执行timer
func (t *Timer) Run()  {
	go func() {
		now := millisecondUnix()

		if t.callTime > now {
			time.Sleep(time.Duration(t.callTime-now) * time.Millisecond)
		}

		t.task.Call()
	}()
}