package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/snorristurluson/exsim_commands"
	"testing"
)

func TestCanCreatePlayer(t *testing.T) {
	player := NewPlayer(1, nil)
	if player.UserId != 1 {
		t.Fail()
	}
}

func TestNewlyCreatedPlayerHasNoShip(t *testing.T) {
	player := NewPlayer(1, nil)
	ship := player.GetShip()
	if ship != nil {
		t.Errorf("Ship is not null")
	}
}

func TestPlayerHandlesSetTargetLocation(t *testing.T) {
	loc := exsim_commands.Vector3{X: 10, Y: 20, Z: 30}
	input := fmt.Sprintf(`{"command": "settargetlocation", "params": { "location": {"x": %v, "y": %v, "z": %v}}}`+"\n", loc.X, loc.Y, loc.Z)
	inputBuffer := bytes.NewBufferString(input)
	outputBuffer := new(bytes.Buffer)
	conn := bufio.NewReadWriter(bufio.NewReader(inputBuffer), bufio.NewWriter(outputBuffer))
	player := NewPlayer(1, conn)
	ship := NewShip(1)
	player.SetShip(ship)
	player.HandleCommand()
	conn.Writer.Flush()
	result := outputBuffer.String()
	if result != "" {
		t.Errorf("Expected '', got '%v'", result)
	}
}
