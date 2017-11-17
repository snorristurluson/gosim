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
	expected := `{"Command":"addship","Params":{"Owner":1,"TypeId":0,"Position":{"X":0,"Y":0,"Z":0}}}` + "\n"
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

	expected := `{"Command":"removeship","Params":{"Owner":1}}` + "\n"
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

	expected := `{"Command":"stepsimulation","Params":{"Timestep":1}}` + "\n"
	expected += `{"Command":"getstate","Params":null}` + "\n"
	if result != expected {
		t.Errorf("Expected '%v', got '%v'", expected, result)
	}
}