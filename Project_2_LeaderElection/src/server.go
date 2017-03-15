package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"

	"shared" //Path to the package contains shared struct
)

func registerArith(server *rpc.Server, arith shared.Arith) {
	// registers Arith interface by name of `Arithmetic`.
	// If you want this name to be same as the type name, you
	// can use server.Register instead.
	server.RegisterName("Arithmetic", arith)
}

func main() {

	//Creating an instance of struct which implement Arith interface
	arith := new(Arith)

	// Register a new rpc server (In most cases, you will use default server only)
	// And register struct we created above by name "Arith"
	// The wrapper method here ensures that only structs which implement Arith interface
	// are allowed to register themselves.
	server := rpc.NewServer()
	registerArith(server, arith)

	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// This statement links rpc server to the socket, and allows rpc server to accept
	// rpc request coming from that socket.
	server.Accept(l)
}
