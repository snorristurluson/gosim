package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"sync"
)

type Login struct {
	User int64
}

type GlobalState struct {
	solarsystems map[string]*Solarsystem
	solarsystemsMutex sync.Mutex
}

func main() {
	globalState := new(GlobalState)
	go globalState.listenForConnections()

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func (gs *GlobalState) listenForConnections() {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		fmt.Printf("Error in Listen: %v\n", err)
		panic(err)
	}
	for {
		fmt.Printf("Waiting for connection on port 4040\n")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error in Accept: %v\n", err)
			continue
		}
		go gs.handleLogin(conn)
	}
}

func (gs *GlobalState) handleLogin(conn net.Conn) {
	var login Login
	for {
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(&login)
		if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			continue
		}

		break
	}

	fmt.Printf("User %v logged in\n", login.User)
	go gs.addPlayer(login.User, conn)
}

func (gs *GlobalState) addPlayer(user int64, conn net.Conn) {
	fmt.Printf("Adding player %v\n", user)
	player := NewPlayer(user, conn)
	ss, err := gs.findSolarsystemForPlayer(player)
	if err != nil {
		fmt.Printf("Can't connect to solar system: %v", err)
		conn.Close()
		return
	}
	ss.AddShip(player.GetShip())
	player.Loop()
}

func (gs *GlobalState) findSolarsystemForPlayer(player *Player) (*Solarsystem, error) {
	gs.solarsystemsMutex.Lock()
	defer gs.solarsystemsMutex.Unlock()

	// Todo: look up solar system where player is located
	ss, ok := gs.solarsystems["ex1"]
	if !ok {
		ss = NewSolarsystem("ex1")
		conn, err := net.Dial("tcp", ":4041")
		if err != nil {
			return nil, err
		}
		ss.SetConnection(conn)
		go ss.Loop()
	}
	return ss, nil
}

func (gs *GlobalState) tickSolarsystems() {

}