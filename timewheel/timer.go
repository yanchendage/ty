package timewheel

import (
	"container/list"
	"sync/atomic"
	"unsafe"
)

//Timer表示单个事件。当Timer超时时，给定的任务将被执行。
type Timer struct {
	expiration int64 // 以毫秒为单位
	task       func() //任务

	b unsafe.Pointer // 所属bucket的指针

	element *list.Element // bucket中timers双向链表中的元素
}
//获取定时器所属的bucket
func (t *Timer) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.b))
}
//设置定时器所属的bucket
func (t *Timer) setBucket(b *bucket) {
	atomic.StorePointer(&t.b, unsafe.Pointer(b))
}

// 阻止定时器启动
func (t *Timer) Stop() bool {
	stopped := false
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		//从bucket（时间格）中移除定时器
		stopped = b.Remove(t)
	}
	return stopped
}

