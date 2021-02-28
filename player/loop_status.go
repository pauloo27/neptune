package player

type LoopStatus int

const (
	LOOP_NONE  = LoopStatus(0)
	LOOP_TRACK = LoopStatus(1)
	LOOP_QUEUE = LoopStatus(2)
)
