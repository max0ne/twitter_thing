package vr

// Command - -
type Command struct {
	Kind  string
	Value interface{}
}

// CommandReply reply of PushCommand
type CommandReply struct {
	Index int
}
