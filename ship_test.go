package main

import (
	"bufio"
	"bytes"
	"github.com/snorristurluson/exsim_commands"
	"testing"
)

func TestCanCreateShip(t *testing.T) {
	ship := NewShip(1)
	if ship.Owner != 1 {
		t.Fail()
	}
}

func TestSetTargetLocationCreatesCommand(t *testing.T) {
	ship := NewShip(1)
	loc := exsim_commands.Vector3{X: 10, Y: 20, Z: 30}
	ship.SetTargetLocation(loc)

	cmds := ship.GetCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected one command, got %v", len(cmds))
	}

	cmd := cmds[0]
	if cmd.Command != "setshiptargetlocation" {
		t.Errorf("Expected 'setshiptargetlocation', got %v", cmd.Command)
	}

	cmds = ship.GetCommands()
	if len(cmds) != 0 {
		t.Errorf("Expected no command, got %v", len(cmds))
	}
}

func TestShipSendState(t *testing.T) {
	inputBuffer := new(bytes.Buffer)
	outputBuffer := new(bytes.Buffer)
	conn := bufio.NewReadWriter(bufio.NewReader(inputBuffer), bufio.NewWriter(outputBuffer))

	ship := NewShip(1)
	ship.SetConnection(conn)

	state := exsim_commands.NewState()
	state.Ships["ship_1"] = exsim_commands.ShipData{
		Owner:    1,
		TypeId:   1,
		Position: exsim_commands.Vector3{X: 1, Y: 2, Z: 3},
		InRange:  []int64{2, 3},
	}
	state.Ships["ship_2"] = exsim_commands.ShipData{
		Owner:    2,
		TypeId:   2,
		Position: exsim_commands.Vector3{X: 4, Y: 5, Z: 6},
		InRange:  []int64{1, 3},
	}
	state.Ships["ship_3"] = exsim_commands.ShipData{
		Owner:    3,
		TypeId:   3,
		Position: exsim_commands.Vector3{X: 7, Y: 8, Z: 9},
		InRange:  []int64{1, 2},
	}

	ship.SendState(state)

	conn.Writer.Flush()
	result := outputBuffer.String()

	expected := `{"ships":{"ship_1":{"owner":1,"typeid":1,"position":{"x":1,"y":2,"z":3},"inrange":[2,3]},"ship_2":{"owner":2,"typeid":2,"position":{"x":4,"y":5,"z":6},"inrange":[1,3]},"ship_3":{"owner":3,"typeid":3,"position":{"x":7,"y":8,"z":9},"inrange":[1,2]}}}` + "\n"
	if result != expected {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}
