package server

import (
	"errors"
	"io"
	"sync"
)

type Pool struct {
	m			sync.Mutex
	resources 	chan io.Closer
	factory 	func()(io.Closer, error)
	closed 		bool
	size		uint32
}

func NewPool(factory func()(io.Closer, error), size uint32) (*Pool, error) {
	if size < 0 {
		return nil, errors.New("【Pool】size value too small")
	}
	
	return &Pool{
		resources : make(chan io.Closer, size),
		factory : factory,
		size : size,
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error)  {
	select {
		case r, ok := <- p.resources:
			if ! ok {
				return  nil, errors.New("【Pool】has been closed")
			}
			return  r, nil
		default:
			return p.factory()
	}
}


func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
		case p.resources <- r:
		default:
			r.Close()
	}
}

func (p *Pool) Close()  {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	close(p.resources)
	for r := range p.resources{
		r.Close()
	}
}
