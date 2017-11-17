package main

import (
	"io"
	"encoding/json"
	"sync"
	"time"
	"fmt"
)

type Solarsystem struct {
	Name  string
	Ships map[int64]*Ship
	Connection io.ReadWriter

	shipsMutex sync.Mutex
	isTicking bool
	isTickingMutex sync.Mutex
}

func NewSolarsystem(name string) *Solarsystem {
	return &Solarsystem{
		Name:  name,
		Ships: make(map[int64]*Ship),
		isTicking: false,
	}
}

func (ss *Solarsystem) SetConnection(conn io.ReadWriter) {
	ss.Connection = conn
	ss.sendCommand(NewSetMainCommand())
}

func (ss *Solarsystem) sendCommand(cmd *Command) (*CommandResult, error) {
	fmt.Printf("Sending command: %v\n", cmd.Command)

	decoder := json.NewDecoder(ss.Connection)
	encoder := json.NewEncoder(ss.Connection)
	encoder.Encode(cmd)

	var result CommandResult
	err := decoder.Decode(&result)
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

	ss.sendCommand(cmd)
}

func (ss *Solarsystem) RemoveShip(ship *Ship) {
	ss.shipsMutex.Lock()
	delete(ss.Ships, ship.Owner)
	ss.shipsMutex.Unlock()

	cmd := NewRemoveShipCommand(ship.Owner)
	ss.sendCommand(cmd)
}

func (ss *Solarsystem) GetShipCount() int {
	return len(ss.Ships)
}

func (ss *Solarsystem) Tick(dt int) error {
	for _, ship := range(ss.Ships) {
		cmds := ship.GetCommands()
		for _, cmd := range(cmds) {
			ss.sendCommand(cmd)
		}
	}

	cmd := NewStepSimulationCommand(float32(dt) / 1000.0)
	ss.sendCommand(cmd)

	cmd = NewGetStateCommand()
	ss.sendCommand(cmd)

	return nil
}

func (ss *Solarsystem) Loop() error {
	tickInterval := 1000 * time.Millisecond
	tickCounter := int64(0)
	for {
		start := time.Now()
		tickCounter += 1

		fmt.Printf("%v: Starting tick %v\n", ss.Name, tickCounter)
		err := ss.Tick(int(tickInterval / time.Millisecond))
		if err != nil {
			fmt.Printf("Error ticking solar system %v: %v", ss.Name, err)
			return err
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
	ss.isTickingMutex.Lock()
	defer ss.isTickingMutex.Unlock()
	ss.isTicking = false

	return nil
}

func (ss *Solarsystem) Start() {
	ss.isTickingMutex.Lock()
	defer ss.isTickingMutex.Unlock()

	if !ss.isTicking {
		ss.isTicking = true
		go ss.Loop()
	}
}