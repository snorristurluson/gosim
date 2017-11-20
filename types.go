package main

import "encoding/json"

type CommandReceived struct {
	Command string `json:"command"`
	Params json.RawMessage `json:"params"`
}

type ParamsReceivedSetTargetLocation struct {
	Location Vector3 `json:"location"`
}

type Command struct {
	Command string `json:"command"`
	Params interface{} `json:"params"`
}

func NewCommand(name string, params interface{}) (*Command){
	return &Command{
		Command: name,
		Params: params,
	}
}

type CommandResult struct {
	Result string `json:"result"`
	State  json.RawMessage `json:"state"`
}

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ShipData struct {
	Owner    int64 `json:"owner"`
	TypeId     int64 `json:"typeid"`
	Position Vector3 `json:"position"`
	InRange  []int64 `json:"inrange"`
}

type State struct {
	Ships map[string]ShipData `json:"ships"`
}

func NewState() (*State) {
	return &State {
		Ships: make(map[string]ShipData),
	}
}

type ParamsAddShip struct {
	Owner int64 `json:"owner"`
	TypeId int64 `json:"typeid"`
	Position Vector3 `json:"position"`
}

func NewAddShipCommand(owner int64, typeid int64, pos Vector3) (*Command){
	params := &ParamsAddShip{
		Owner: owner,
		TypeId: typeid,
		Position: pos,
	}
	return NewCommand("addship", params)
}

type ParamsRemoveShip struct {
	Owner int64 `json:"owner"`
}

func NewRemoveShipCommand(owner int64) (*Command) {
	params := &ParamsRemoveShip{
		Owner: owner,
	}
	return NewCommand("removeship", params)
}

type ParamsSetShipTargetLocation struct {
	ShipId int64 `json:"shipid"`
	Location Vector3 `json:"location"`
}

func NewSetShipTargetLocationCommand(shipid int64, loc Vector3) (*Command) {
	params := &ParamsSetShipTargetLocation{
		ShipId: shipid,
		Location: loc,
	}
	return NewCommand("setshiptargetlocation", params)
}

type ParamsStepSimulation struct {
	Timestep float32 `json:"timestep"`
}

func NewStepSimulationCommand(timestep float32) (*Command) {
	params := &ParamsStepSimulation{
		Timestep: timestep,
	}
	return NewCommand("stepsimulation", params)
}

func NewGetStateCommand() (*Command) {
	return NewCommand("getstate", nil)
}

func NewSetMainCommand() (*Command) {
	return NewCommand("setmain", nil)
}