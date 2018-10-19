package varmint

type Status int

const (
	// The Runner's function is currently executing
	Running Status = iota
	// The Runner is waiting for some condition before executing
	Waiting Status = iota
	// The runner is currently paused and will not run until restarted
	Paused Status = iota
	// The runner is cancelled and will not run again
	Cancelled Status = iota
)

func (s Status) String() string {
	switch s {
	case Running:
		return "Running"
	case Waiting:
		return "Waiting"
	case Paused:
		return "Paused"
	case Cancelled:
		return "Cancelled"
	default:
		return "UNKNOWN"
	}
}
