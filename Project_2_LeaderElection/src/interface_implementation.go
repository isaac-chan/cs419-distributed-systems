package main

import (
     "shared" //Path to the package contains shared struct
)

// Every method that we want to export must have
// (1) the method has two arguments, both exported (or builtin) types
// (2) the method's second argument is a pointer
// (3) the method has return type error

type Arith int

func (t *Arith) Multiply(args *shared.Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *shared.Args, quo *shared.Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
