package timewheel

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

// 优先级队列实现的开始。
// The start of PriorityQueue implementation.
// Borrowed from https://github.com/nsqio/nsq/blob/master/internal/pqueue/pqueue.go

//队列中的元素
type item struct {
	Value    interface{}
	Priority int64 //优先权
	Index    int
}

//这是由最小堆实现的优先队列
type priorityQueue []*item

func newPriorityQueue(capacity int) priorityQueue {
	return make(priorityQueue, 0, capacity)
}

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	//动态扩容
	if n+1 > c {
		npq := make(priorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}

	*pq = (*pq)[0 : n+1]
	item := x.(*item)
	item.Index = n
	(*pq)[n] = item
}

func (pq *priorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	//动态缩容
	if n < (c/2) && c > 25 {
		npq := make(priorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}

	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

// 与最小堆的堆顶比较
// 如果当前时间小于最小堆的堆顶，说明堆里所有的元素均没有到过期时间
// 如果当前时间大于最小堆的堆顶，移出堆顶，并重新排序

func (pq *priorityQueue) PeekAndShift(max int64) (*item, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	//最小堆的顶大于max
	if item.Priority > max {
		return nil, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}

// 优先级队列实现的结束。

//延迟队列
type DelayQueue struct {
	C chan interface{} // 有元素过期时的通知

	mu sync.Mutex	// 互斥锁
	pq priorityQueue // 优先队列

	sleeping int32 // 已休眠
	wakeupC  chan struct{} // 唤醒队列的通知
}


func New(size int) *DelayQueue {

	return &DelayQueue{
		C:       make(chan interface{}), // 无缓冲管道
		pq:      newPriorityQueue(size), // 优先队列
		wakeupC: make(chan struct{}), // 无缓冲管道saw
	}
}

//添加元素到队列
func (dq *DelayQueue) Offer(elem interface{}, expiration int64) {
	item := &item{Value: elem, Priority: expiration} //过期时间作为优先级，过期时间越小的优先级越高

	dq.mu.Lock()
	heap.Push(&dq.pq, item) // 将元素放到队尾，并递归与父节点做比较
	index := item.Index
	dq.mu.Unlock()

	if index == 0 {
		// 如果延迟队列为休眠状态，唤醒他
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			// 唤醒可能会发生阻塞
			dq.wakeupC <- struct{}{}
		}
	}
}

// Poll启动一个无限循环，在这个循环中它不断地等待一个元素过期，然后将过期的元素发送到通道C。
func (dq *DelayQueue) Poll(exitC chan struct{}, nowF func() int64) {
	for {
		now := nowF() // nowF是timeToMs(time.Now().UTC())，当前时间的毫秒时间戳

		dq.mu.Lock()
		item, delta := dq.pq.PeekAndShift(now) //与最小堆的堆顶比较

		if item == nil {
			//没有要过期的定时器，	将延迟队列设置为休眠
			//为什么要用atomic原子函数，是为了防止Offer 和 Poll出现竞争
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item == nil {
			if delta == 0 {
				// 说明延迟队列中已经没有timer，因此等待新的timer添加时wake up通知，或者等待退出通知
				select {
				case <-dq.wakeupC:
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 {
				// 说明延迟队列中存在未过期的定时器
				select {
				case <-dq.wakeupC:
					// 当前定时器已经是休眠状态，如果添加了一个比延迟队列中最早过期的定时器更早的定时器,延迟队列被唤醒
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// timer.After添加了一个相对时间定时器,并等待到期

					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						//防止被阻塞
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		select {
		case dq.C <- item.Value:
		case <-exitC:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
}