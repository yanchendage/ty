package timewheel

import (
	"sync"
	"time"
)

//将x截取到m的最大整数备份
//x = 5, m = 3 , 5 - 5%3 = 3
//x = 100, m = 3, 100 - 100%3 = 99
func truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}

//返回1970-1-1至今的毫秒数
func timeToMs(t time.Time) int64 {
	//return 	return time.Now().UnixNano()  / 1e6
	return t.UnixNano() / int64(time.Millisecond)
}

// msToTime returns the UTC time corresponding to the given Unix time,
// t milliseconds since January 1, 1970 UTC.
func msToTime(t int64) time.Time {
	return time.Unix(0, t*int64(time.Millisecond)).UTC()
}

type waitGroupWrapper struct {
	sync.WaitGroup
}

//包裹
func (w *waitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
