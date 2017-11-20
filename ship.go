package main

import (
	"io"
	"fmt"
	"encoding/json"
)

type Ship struct {
	Owner int64
	TypeId int64
	Position Vector3
	InRange []int64
	commands []*Command
	connection io.ReadWriter
}

func NewShip(owner int64) *Ship {
	return &Ship{
		Owner: owner,
	}
}

func (ship *Ship) SetConnection(conn io.ReadWriter) {
	ship.connection = conn
}

func (ship *Ship) GetCommands() ([]*Command) {
	commands := ship.commands
	ship.commands = make([]*Command, 0)
	return commands
}

func (ship *Ship) SetTargetLocation(loc Vector3) {
	cmd := NewSetShipTargetLocationCommand(ship.Owner, loc)
	ship.commands = append(ship.commands, cmd)
}

// This should only be called before the ship is added to the solar system.
func (ship *Ship) SetPosition(pos Vector3) {
	ship.Position = pos
}

// Send state to player
func (ship *Ship) SendState(state *State) {
	me := state.Ships[fmt.Sprintf("ship_%v", ship.Owner)]
	ship.Position = me.Position
	ship.InRange = me.InRange

	// Gather up my view of the world
	myWorldState := NewState()

	myWorldState.Ships[fmt.Sprintf("ship_%v", me.Owner)] = me
	for _, v := range me.InRange {
		key := fmt.Sprintf("ship_%v", v)
		other := state.Ships[key]
		myWorldState.Ships[key] = other
	}

	encoder := json.NewEncoder(ship.connection)
	encoder.Encode(myWorldState)
}
