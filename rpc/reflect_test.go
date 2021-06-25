package rpc

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestReflect(t *testing.T) {
	var wg sync.WaitGroup

	typ := reflect.TypeOf(&wg)

	for i:=0; i < typ.NumMethod(); i++{
		method := typ.Method(i)

		args:= make([]string,0, method.Type.NumIn())
		returns := make([]string,0, method.Type.NumOut())

		for j := 1; j < method.Type.NumIn(); j++ {
			args = append(args, method.Type.In(j).Name())
		}

		for j := 0; j < method.Type.NumOut(); j++ {
			returns = append(returns, method.Type.Out(j).Name())
		}

		fmt.Println(strings.Join(args,","))
	}
}


type Foo int

type Args struct{ Num1, Num2 int }
type Fuck struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, r *int) error {
	*r = args.Num1 + args.Num2
	return nil
}

// it's not a exported Method
func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func TestNewService(t *testing.T)  {
	var foo Foo

	s := newService(&foo)

	methodType := s.method["Sum"]

	if methodType == nil {
		log.Println("")
	}
}

func TestMethodTypeCall(t *testing.T)  {

	var  foo Foo
	s := newService(&foo)

	methodType := s.method["Sum"]

	arg := methodType.newArgv()
	reply := methodType.newReplyv()

	arg.Set(reflect.ValueOf(Args{1,1}))

	s.call(methodType,arg,reply)

	log.Println(*reply.Interface().(*int))

}

