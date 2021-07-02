package main

import "strings"

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