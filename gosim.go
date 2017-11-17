package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"math/rand"
)

type Login struct {
	User int64
}

type GlobalState struct {
	solarsystems map[string]*Solarsystem
	solarsystemsMutex sync.Mutex
}

func NewGlobalState() (*GlobalState){
	return &GlobalState{
		solarsystems: make(map[string]*Solarsystem),
	}
}

func main() {
	globalState := NewGlobalState()
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

	// Todo: ship position should be looked up
	ship := player.GetShip()
	ship.SetPosition(Vector3{X:rand.Float64()*5000.0-2500.0, Y:rand.Float64()*5000.0-2500.0, Z: 0})
	ss.AddShip(player.GetShip())
	player.Loop()
	ss.RemoveShip(player.GetShip())
}

func (gs *GlobalState) findSolarsystemForPlayer(player *Player) (*Solarsystem, error) {
	gs.solarsystemsMutex.Lock()
	defer gs.solarsystemsMutex.Unlock()

	// Todo: look up solar system where player is located
	name := "ex1"
	ss, ok := gs.solarsystems[name]
	if !ok {
		ss = NewSolarsystem(name)
		conn, err := net.Dial("tcp", ":4041")
		if err != nil {
			return nil, err
		}
		gs.solarsystems[name] = ss
		ss.SetConnection(conn)
	}
	ss.Start()
	return ss, nil
}
