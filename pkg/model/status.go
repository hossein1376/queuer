package model

// Status represent the [Order] processing status. Its zero is valid and
// indicates the [Pending] status.
type Status int16

const (
	Pending Status = iota
	Processed
	Failed
)
