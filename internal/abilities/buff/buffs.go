package buff

import "time"

type Buff struct {
	Duration time.Duration
	ExpireAt time.Time
	Charges  int
	Enable   func(*Status)
	Disable  func(*Status)
}

type Status struct {
	Invulnerable int
}

func Invulnerable(dur time.Duration) Buff {
	return Buff{
		Duration: dur,
		Enable: func(s *Status) {
			s.Invulnerable++
		},
		Disable: func(s *Status) {
			s.Invulnerable--
		},
	}
}
