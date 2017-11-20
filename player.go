package main

import (
	"fmt"
	"io"
	"encoding/json"
	"errors"
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
	var cmd CommandReceived
	err := decoder.Decode(&cmd)
	if err != nil {
		fmt.Printf("Error reading command: %v\n", err)
		player.Connection.Write([]byte(`{"result": "error"}` + "\n"))
		return err
	}

	if cmd.Command == "settargetlocation" {
		var params ParamsReceivedSetTargetLocation
		err := json.Unmarshal(cmd.Params, &params)
		if err != nil {
			fmt.Printf("Error unmarshaling params: %v", err)
		}
		player.Ship.SetTargetLocation(params.Location)
		player.Connection.Write([]byte(`{"result": "ok"}` + "\n"))
	} else {
		fmt.Printf("Error reading command: %v\n", err)
		player.Connection.Write([]byte(`{"result": "error"}` + "\n"))
		return errors.New("Unknown command")
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

