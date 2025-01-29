package model

// Priority dictates the priority level of [Order] processing. Its zero
// value is valid and have the lowest priority, as in [Normal].
type Priority int16

const (
	Normal Priority = iota
	High
)
