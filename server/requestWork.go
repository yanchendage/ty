package server

type requestTask struct {
	request IRequest
	router IRouter
}

func (rt *requestTask) task()  {
	rt.router.PreHandle(rt.request)
	rt.router.Handle(rt.request)
	rt.router.PostHandle(rt.request)
}
