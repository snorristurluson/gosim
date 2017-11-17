package main

type Ship struct {
	Owner int64
	TypeId int64
	Position Vector3
	commands []*Command
}

func NewShip(owner int64) *Ship {
	return &Ship{
		Owner: owner,
	}
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