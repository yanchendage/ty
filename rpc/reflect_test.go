package rpc

import (
	"bytes"
	"encoding/gob"
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

type Param struct{ Num1, Num2 int }

func (f Foo) Sum(args Param, r *int) error {
	*r = args.Num1 + args.Num2
	return nil
}
func (f Foo) Echo(args []string, r *string) error {
	*r = strings.Join(args,",")
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

	arg.Set(reflect.ValueOf(Param{1,1}))

	s.call(methodType,arg,reply)

	log.Println(*reply.Interface().(*int))
}

func TestReflect2(t *testing.T)  {

	var  foo Foo
	s := newService(&foo)

	methodType := s.method["Sum"]

	arg := methodType.newArgv()
	reply := methodType.newReplyv()

	arg.Set(reflect.ValueOf(Param{1,1}))

	s.call(methodType,arg,reply)

	log.Println(*reply.Interface().(*int))
}


func TestStruct(t *testing.T)  {
	type People struct {
		Name string
		Age int
		Like []string
		Nvyou map[string]string
	}

	m:=make(map[string]string)
	m["S"]="hello"
	m["SS"]="jiangzhou"


	p := People{"yanchen",29, []string{"changge","tiaowu1"}, m}
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	//gob.Register(map[string]string{})
	err := enc.Encode(&p)
	if err !=nil {
		fmt.Println(err)
	}

	var p2 People
	dec := gob.NewDecoder(buf)

	dec.Decode(&p2)

	fmt.Println(p)
	fmt.Println(p2)
}


func TestStructV2(t *testing.T)  {
	type People struct {
		Name string
		Age int
		Like []string
		Nvyou interface{}
	}

	m:=make(map[string]string)
	m["S"]="hello"
	m["SS"]="jiangzhou"




	p := People{"yanchen",29, []string{"changge","tiaowu1"}, m}
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	//gob.Register(map[string]string{})
	err := enc.Encode(&p.Name)
	err = enc.Encode(&p.Age)
	err = enc.Encode(&p.Nvyou)

	if err !=nil {
		fmt.Println(err)
	}

	var p2 People
	dec := gob.NewDecoder(buf)

	dec.Decode(&p2.Name)
	dec.Decode(&p2.Age)
	dec.Decode(&p2.Nvyou)

	fmt.Println(p)
	fmt.Println(p2.Age)
	fmt.Println(p2.Nvyou)
}

func TestStructV3(t *testing.T)  {
	type People struct {
		Name string
		Age int
		Like []string
		Nvyou interface{}
	}

	m:=make(map[string]string)
	m["S"]="hello"
	m["SS"]="jiangzhou"


	p := People{"yanchen",29, []string{"changge","tiaowu1"}, m}
	buf := new(bytes.Buffer)

	enc := gob.NewEncoder(buf)
	//gob.Register(map[string]string{})
	err := enc.Encode(&p.Name)
	err = enc.Encode(&p.Age)
	err = enc.Encode(&p.Nvyou)

	if err !=nil {
		fmt.Println(err)
	}

	var p2 People
	dec := gob.NewDecoder(buf)

	dec.Decode(&p2.Name)
	dec.Decode(&p2.Age)
	dec.Decode(&p2.Nvyou)

	fmt.Println(p)
	fmt.Println(p2.Age)
	fmt.Println(p2.Nvyou)
}

func TestR(t *testing.T)  {
	type R struct {
		R1 reflect.Value
		R2 reflect.Value
	}

	var R3 int

	e := &R{reflect.ValueOf(123), reflect.ValueOf([]string{"yanchen"})}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.EncodeValue(e.R1)

	fmt.Println(buf)
	dec := gob.NewDecoder(buf)

	//d2 := &R{}
	err := dec.DecodeValue(reflect.ValueOf(&R3))
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.ValueOf(R3))
	//dec.DecodeValue(d.R2)
}

func TestR2(t *testing.T)  {
	type R struct {
		R1 reflect.Value
		R2 reflect.Value
	}
	type FFF struct {
		Num1 int
		Num2 int
	}
	f := FFF{1,2}

	e := &R{reflect.ValueOf(f),reflect.ValueOf([]string{"yanchen"}) }

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.EncodeValue(e.R1)
	enc.EncodeValue(e.R2)

	dec := gob.NewDecoder(buf)

	var a FFF
	var argv reflect.Value
	argType := reflect.TypeOf(a)

	var b []string
	var brgv reflect.Value
	brgType := reflect.TypeOf(b)

	if argType.Kind() == reflect.Ptr {
		argv = reflect.New(argType.Elem())
	} else {
		argv = reflect.New(argType).Elem()
	}
	argvi := argv.Interface()
	if argv.Type().Kind() != reflect.Ptr {
		argvi = argv.Addr().Interface()
	}

	if brgType.Kind() == reflect.Ptr {
		brgv = reflect.New(brgType.Elem())
	} else {
		brgv = reflect.New(brgType).Elem()
	}
	brgvi := brgv.Interface()
	if brgv.Type().Kind() != reflect.Ptr {
		brgvi = brgv.Addr().Interface()
	}




	//
	//argv = reflect.New(argType).Elem()
	//argvi := argv.Addr().Interface()
	//fmt.Println(argvi)

	//argv = reflect.New(argv2Type).Elem()
	//argvi2 := argv2.Interface()
	//fmt.Println(argvi)


	//if m.ArgType.Kind() == reflect.Ptr {
	//	argv = reflect.New(m.ArgType.Elem())
	//} else {
	//	argv = reflect.New(m.ArgType).Elem()
	//}

	//err := dec.DecodeValue(reflect.ValueOf(argvi2))
	err := dec.DecodeValue(reflect.ValueOf(argvi))
	err = dec.DecodeValue(reflect.ValueOf(brgvi))

	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.ValueOf(argvi).Elem())
	fmt.Println(reflect.ValueOf(brgvi).Elem())
}
