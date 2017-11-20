package main

import (
	"io"
	"encoding/json"
	"sync"
	"time"
	"fmt"
	"net"
)

type Solarsystem struct {
	Name  string
	Ships map[int64]*Ship

	connection io.ReadWriter
	encoder    *json.Encoder
	decoder    *json.Decoder
	sendQueue chan *Command

	shipsMutex     sync.Mutex
	isTicking      bool
	isTickingMutex sync.Mutex
}

func NewSolarsystem(name string) *Solarsystem {
	return &Solarsystem{
		Name:      name,
		Ships:     make(map[int64]*Ship),
		isTicking: false,
		sendQueue: make(chan *Command, 10),
	}
}

func (ss *Solarsystem) SetConnection(conn io.ReadWriter) {
	ss.connection = conn
	ss.decoder = json.NewDecoder(ss.connection)
	ss.encoder = json.NewEncoder(ss.connection)
	ss.sendCommand(NewSetMainCommand())
}

func (ss *Solarsystem) sendCommand(cmd *Command) (*CommandResult, error) {
	fmt.Printf("Sending command: %v\n", cmd.Command)

	ss.encoder.Encode(cmd)

	var result CommandResult
	err := ss.decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ss *Solarsystem) AddShip(ship *Ship) {
	ss.shipsMutex.Lock()
	ss.Ships[ship.Owner] = ship
	ss.shipsMutex.Unlock()

	cmd := NewAddShipCommand(ship.Owner, ship.TypeId, ship.Position)

	ss.sendQueue <- cmd
}

func (ss *Solarsystem) RemoveShip(ship *Ship) {
	ss.shipsMutex.Lock()
	delete(ss.Ships, ship.Owner)
	ss.shipsMutex.Unlock()

	cmd := NewRemoveShipCommand(ship.Owner)
	ss.sendQueue <- cmd
}

func (ss *Solarsystem) GetShipCount() int {
	return len(ss.Ships)
}

func (ss *Solarsystem) Tick(dt int) error {
	ss.HandleQueuedCommands()

	for _, ship := range (ss.Ships) {
		cmds := ship.GetCommands()
		for _, cmd := range (cmds) {
			ss.sendCommand(cmd)
		}
	}

	cmd := NewStepSimulationCommand(float32(dt) / 1000.0)
	result, err := ss.sendCommand(cmd)
	if err != nil {
		fmt.Printf("Error in stepsimulation: %v\n", err)
		return err
	}

	cmd = NewGetStateCommand()
	result, err = ss.sendCommand(cmd)
	if err != nil {
		fmt.Printf("Error in getstate: %v\n", err)
		return err
	}

	var state State
	json.Unmarshal(result.State, &state)

	ss.shipsMutex.Lock()
	defer ss.shipsMutex.Unlock()

	var wg sync.WaitGroup
	for _, ship := range (ss.Ships) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ship.SendState(&state)
		}()
	}

	ss.HandleQueuedCommands()

	return nil
}

func (ss *Solarsystem) HandleQueuedCommands() {
	for {
		select {
		case cmd := <-ss.sendQueue:
			ss.sendCommand(cmd)
		default:
			return
		}
	}
}

func (ss *Solarsystem) Loop() error {
	addr, err := net.ResolveTCPAddr("tcp", ":4041")
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	ss.SetConnection(conn)

	err = nil
	tickInterval := 1000 * time.Millisecond
	tickCounter := int64(0)
	for {
		start := time.Now()
		tickCounter += 1

		fmt.Printf("%v: Starting tick %v\n", ss.Name, tickCounter)
		err := ss.Tick(int(tickInterval / time.Millisecond))
		if err != nil {
			fmt.Printf("Error ticking solar system %v: %v\n", ss.Name, err)
			break
		}
		if len(ss.Ships) == 0 {
			fmt.Printf("System is empty, stopping loop\n")
			break
		}
		tickDuration := time.Since(start)
		sleepTime := tickInterval - tickDuration
		if sleepTime > 0 {
			fmt.Printf("Tick duration: %v - Sleeping for %v milliseconds\n", int(tickDuration/time.Millisecond), int(sleepTime/time.Millisecond))
			time.Sleep(sleepTime)
		}
	}
	conn.Close()
	ss.connection = nil

	ss.isTickingMutex.Lock()
	defer ss.isTickingMutex.Unlock()
	ss.isTicking = false

	return err
}

func (ss *Solarsystem) Start() {
	ss.isTickingMutex.Lock()
	defer ss.isTickingMutex.Unlock()

	if !ss.isTicking {
		ss.isTicking = true
		go ss.Loop()
	}
}
