package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/snorristurluson/exsim_commands"
	"io"
)

type Player struct {
	UserId     int64
	Connection io.ReadWriter
	Ship       *Ship
}

func NewPlayer(user int64, conn io.ReadWriter) *Player {
	return &Player{
		UserId:     user,
		Connection: conn,
	}
}

func (player *Player) SetShip(ship *Ship) {
	player.Ship = ship
}

func (player *Player) GetShip() *Ship {
	return player.Ship
}

func (player *Player) HandleCommand() error {
	decoder := json.NewDecoder(player.Connection)
	var cmd exsim_commands.CommandReceived
	err := decoder.Decode(&cmd)
	if err != nil {
		fmt.Printf("Error reading command: %v\n", err)
		player.Connection.Write([]byte(`{"result": "error"}` + "\n"))
		return err
	}

	fmt.Printf("%v\n", cmd.Command)
	if cmd.Command == "settargetlocation" {
		var params exsim_commands.ParamsReceivedSetTargetLocation
		err := json.Unmarshal(cmd.Params, &params)
		if err != nil {
			fmt.Printf("Error unmarshaling params: %v\n", err)
			return err
		}
		player.Ship.SetTargetLocation(params.Location)
	} else if cmd.Command == "setattribute" {
		var params exsim_commands.ParamsReceivedSetAttribute
		err := json.Unmarshal(cmd.Params, &params)
		if err != nil {
			fmt.Printf("Error unmarshaling params: %v\n", err)
			return err
		}
		player.Ship.SetAttribute(params.Attribute, params.Value)
	} else {
		fmt.Printf("Error reading command: %v\n", err)
		return errors.New("unknown command")
	}

	return nil
}

func (player *Player) Loop() {
	for {
		err := player.HandleCommand()
		if err == io.EOF {
			fmt.Printf("Player connection %v lost\n", player.UserId)
			return
		}
	}
}
