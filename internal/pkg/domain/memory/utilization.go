package memory

type Utilization struct {
	Total uint64
	Free  uint64
	Used  uint64
}

func (u *Utilization) Percent() float64 {
	return float64(u.Used) / float64(u.Total)
}
