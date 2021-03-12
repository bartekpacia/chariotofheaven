package main

type MoveState string

const (
	MoveForward  MoveState = "w"
	MoveBackward MoveState = "s"
	MoveStop     MoveState = "q"
)

type TurnState string

const (
	TurnLeft  TurnState = "a"
	TurnRight TurnState = "d"
	TurnStop  TurnState = "z"
)

type Chariot struct {
	moveState MoveState
	turnState TurnState
}
