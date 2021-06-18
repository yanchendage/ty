package server

import (
	"log"
)

type MessageManager struct {
	//msgID 2 router
	maps map[uint32]IRouter
	//number of workerPool
	WorkerPoolNumber uint32

	WorkQueue []chan IWorker
}


func NewMessageManager() *MessageManager {
	return &MessageManager{
		maps:make(map[uint32]IRouter),
		WorkerPoolNumber : Conf["workerPoolNumber"].(uint32),
		WorkQueue: make([]chan IWorker, Conf["workerPoolNumber"].(uint32)), // match the TaskQueue with the WorkerPool
	}
}

func (mm *MessageManager) AddRouter(msgID uint32, router IRouter)  {
	if _, ok := mm.maps[msgID]; !ok{
		mm.maps[msgID] = router
	}
}

func  (mm *MessageManager) GetRouter(msgID uint32) IRouter{
	handler, ok := mm.maps[msgID]

	if !ok {
		log.Println("【Message Manager】 msgId = ", msgID , " is not FOUND!")
		return nil
	}

	return handler
}

func (mm *MessageManager) StartWorkerPool() {
	for i := 0; i< int(Conf["workerPoolNumber"].(uint32)); i++ {
		mm.WorkQueue[i] = make(chan IWorker, Conf["maxTaskQueueLen"].(int))

		go mm.startOneWorkerPool(i, mm.WorkQueue[i])

		log.Println("【Message Manager】work pool id", i, "starting")
	}
}

func (mm *MessageManager) startOneWorkerPool(workerPoolID int, taskQueue chan IWorker)  {
	wp := NewWorkerPool(1)

	for {
		select {
		case w := <-taskQueue :
			wp.AddWork(w)
		}
	}
}

func (mm *MessageManager) AddWorkerToWorkQueue(request IRequest){
	// balance
	workQueueId := request.GetConnection().GetConnID() %  mm.WorkerPoolNumber

	log.Println("【Message Manager】Add ConnID=",
		request.GetConnection().GetConnID(),
		" request msgID=",
		request.GetMsgID(),
		"to workerQueueID=",
		workQueueId)

	handler := mm.GetRouter(request.GetMsgID())

	rt := requestTask{
		request: request,
		router:  handler,
	}

	mm.WorkQueue[workQueueId] <- &rt
}


