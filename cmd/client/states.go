package main

const (
	MovingForward = iota + 1
	MovingBackward
	NotMoving
)

const (
	TurningLeft = iota + 1
	TurningRight
	NotTurning
)

type ChariotState struct {
	MovingState  int
	TurningState int
}

// MapCommandToState maps a command (a single char) to chariot's state.
func (c *ChariotState) MapCommandToState(cmd string) {
	switch cmd {
	// Move commands
	case CmdMoveForward:
		c.MovingState = MovingForward
	case CmdMoveBackward:
		c.MovingState = MovingBackward
	case CmdMoveStop:
		c.MovingState = NotMoving
	// Turn commands
	case CmdTurnLeft:
		c.TurningState = TurningLeft
	case CmdTurnRight:
		c.TurningState = TurningRight
	case CmdTurnStop:
		c.TurningState = NotTurning
	// Stop command
	case CmdStopAll:
		c.MovingState = NotMoving
		c.TurningState = NotTurning
	}
}
