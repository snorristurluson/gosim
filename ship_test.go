package main

import "testing"

func TestCanCreateShip(t *testing.T) {
	ship := NewShip(1)
	if ship.Owner != 1 {
		t.Fail()
	}
}

func TestSetTargetLocationCreatesCommand(t *testing.T) {
	ship := NewShip(1)
	loc := Vector3{X: 10, Y: 20, Z: 30}
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