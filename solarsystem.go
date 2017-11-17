package main

import (
	"io"
	"encoding/json"
	"sync"
	"time"
)

type Solarsystem struct {
	Name  string
	Ships map[int64]*Ship
	Connection io.ReadWriter

	shipsMutex sync.Mutex
}

func NewSolarsystem(name string) *Solarsystem {
	return &Solarsystem{
		Name:  name,
		Ships: make(map[int64]*Ship),
	}
}

func (ss *Solarsystem) SetConnection(conn io.ReadWriter) {
	ss.Connection = conn
}

func (ss *Solarsystem) AddShip(ship *Ship) {
	ss.shipsMutex.Lock()
	defer ss.shipsMutex.Unlock()

	ss.Ships[ship.Owner] = ship

	cmd := NewAddShipCommand(ship.Owner, ship.TypeId, ship.Position)

	encoder := json.NewEncoder(ss.Connection)
	encoder.Encode(cmd)
}

func (ss *Solarsystem) RemoveShip(ship *Ship) {
	ss.shipsMutex.Lock()
	defer ss.shipsMutex.Unlock()

	delete(ss.Ships, ship.Owner)

	cmd := NewRemoveShipCommand(ship.Owner)

	encoder := json.NewEncoder(ss.Connection)
	encoder.Encode(cmd)
}

func (ss *Solarsystem) GetShipCount() int {
	return len(ss.Ships)
}

func (ss *Solarsystem) Tick(dt float32) error {
	var result CommandResult
	decoder := json.NewDecoder(ss.Connection)
	encoder := json.NewEncoder(ss.Connection)

	for _, ship := range(ss.Ships) {
		cmds := ship.GetCommands()
		for cmd := range(cmds) {
			encoder.Encode(cmd)
		}
	}

	cmd := NewStepSimulationCommand(dt)
	encoder.Encode(cmd)
	err := decoder.Decode(&result)
	if err != nil {
		return err
	}

	cmd = NewGetStateCommand()
	encoder.Encode(cmd)
	err = decoder.Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func (ss *Solarsystem) Loop() {
	start := time.Now()
	ss.Tick(0.25)
	tickDuration := time.Since(start)
	sleepTime := 250 - tickDuration * time.Millisecond
	time.Sleep(sleepTime * time.Millisecond)
}