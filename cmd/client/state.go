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

func (cs *Chariot) ExecuteCommand(command string) {
	// Move commands
	switch command {
	case MoveForward:
		cs.MovingState = MovingForward
		return

	case MoveBackward:
		cs.MovingState = MovingBackward
		return

	case MoveStop:
		cs.MovingState = NotMoving
		return
	}

	// Turn commands
	switch command {
	case TurnLeft:
		cs.TurningDirection = Left
		cs.Turning = true
		return

	case TurnRight:
		cs.TurningDirection = Right
		cs.Turning = true
		return

	case TurnStop:
		cs.Turning = false
		return
	}

	if command == StopAll {
		cs.MovingState = NotMoving
		cs.Turning = false
		return
	}
}
