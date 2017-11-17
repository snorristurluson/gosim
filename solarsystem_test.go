package main

import (
	"testing"
	"bytes"
	"bufio"
)

func TestCanCreateSolarsystem(t *testing.T) {
	var ss = NewSolarsystem("ex1")
	if ss.Name != "ex1" {
		t.Fail()
	}
	if ss.GetShipCount() != 0 {
		t.Fail()
	}
}

func TestCanAddShipToSolarsystem(t *testing.T) {
	var ss = NewSolarsystem("ex1")

	inputBuffer := new(bytes.Buffer)
	outputBuffer := new(bytes.Buffer)
	conn := bufio.NewReadWriter(bufio.NewReader(inputBuffer), bufio.NewWriter(outputBuffer))

	ss.SetConnection(conn)

	var ship = NewShip(1)
	ss.AddShip(ship)
	if ss.GetShipCount() != 1 {
		t.Fail()
	}

	conn.Writer.Flush()
	result := outputBuffer.String()
	expected := `{"command":"addship","params":{"owner":1,"typeid":0,"position":{"x":0,"y":0,"z":0}}}` + "\n"
	if result != expected {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}

func TestCanRemoveShipFromSolarsystem(t *testing.T) {
	var ss = NewSolarsystem("ex1")

	inputBuffer := new(bytes.Buffer)
	outputBuffer := new(bytes.Buffer)
	conn := bufio.NewReadWriter(bufio.NewReader(inputBuffer), bufio.NewWriter(outputBuffer))

	ss.SetConnection(conn)

	var ship = NewShip(1)
	ss.AddShip(ship)

	conn.Writer.Flush()
	outputBuffer.Reset()

	ss.RemoveShip(ship)

	conn.Writer.Flush()
	result := outputBuffer.String()

	expected := `{"command":"removeship","params":{"owner":1}}` + "\n"
	if result != expected {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}

	if ss.GetShipCount() != 0 {
		t.Fail()
	}
}

func TestSimpleTick(t *testing.T) {
	var ss = NewSolarsystem("ex1")

	inputBuffer := new(bytes.Buffer)
	outputBuffer := new(bytes.Buffer)
	conn := bufio.NewReadWriter(bufio.NewReader(inputBuffer), bufio.NewWriter(outputBuffer))

	ss.SetConnection(conn)

	var ship = NewShip(1)
	ss.AddShip(ship)

	conn.Writer.Flush()
	outputBuffer.Reset()

	inputBuffer.WriteString(`{"Result":"ok"}` + "\n")
	inputBuffer.WriteString(`{"Result":"ok"}` + "\n")

	err := ss.Tick(1)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	conn.Writer.Flush()
	result := outputBuffer.String()

	expected := `{"command":"stepsimulation","params":{"timestep":1}}` + "\n"
	expected += `{"command":"getstate","params":null}` + "\n"
	if result != expected {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}