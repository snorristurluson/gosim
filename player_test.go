package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestCanCreatePlayer(t *testing.T) {
	player := NewPlayer(1, nil)
	if player.UserId != 1 {
		t.Fail()
	}
}

func TestNewlyCreatedPlayerHasShip(t *testing.T) {
	player := NewPlayer(1, nil)
	ship := player.GetShip()
	if ship == nil {
		t.Log("Ship is null")
		t.Fail()
	}
}

func TestPlayerHandlesSetTargetLocation(t *testing.T) {
	loc := Vector3{X: 10, Y: 20, Z: 30}
	input := fmt.Sprintf(`{"command": "settargetlocation", "params": { "location": {"x": %v, "y": %v, "z": %v}}}`+"\n", loc.X, loc.Y, loc.Z)
	inputBuffer := bytes.NewBufferString(input)
	outputBuffer := new(bytes.Buffer)
	conn := bufio.NewReadWriter(bufio.NewReader(inputBuffer), bufio.NewWriter(outputBuffer))
	player := NewPlayer(1, conn)
	player.HandleCommand()
	conn.Writer.Flush()
	expected := `{"result": "ok"}` + "\n"
	result := outputBuffer.String()
	if result != expected {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}
