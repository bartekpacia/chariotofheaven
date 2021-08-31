package main

type MovingState int

const (
	MovingForward MovingState = iota + 1
	MovingBackward
	NotMoving
)

type TurningDirection int

const (
	Left TurningDirection = iota + 1
	Right
)

type Chariot struct {
	MovingState      MovingState
	TurningDirection TurningDirection
	Turning          bool
}

func (cs *Chariot) InterpretCommand(command string) {
	// Move commands
	switch command {
	case MoveForward:
		cs.MovingState = MovingForward
	case MoveBackward:
		cs.MovingState = MovingBackward
	case MoveStop:
		cs.MovingState = NotMoving
	}

	// Turn commands
	switch command {
	case TurnLeft:
		cs.TurningDirection = Left
		cs.Turning = true
	case TurnRight:
		cs.TurningDirection = Right
		cs.Turning = true
	case TurnStop:
		cs.Turning = false
	}

	// Stop command
	if command == StopAll {
		cs.MovingState = NotMoving
		cs.Turning = false
	}
}
