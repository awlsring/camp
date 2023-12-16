package cpu

type Core struct {
	Id      int
	Threads uint32
}

func NewCore(id int, threads uint32) *Core {
	return &Core{
		Id:      id,
		Threads: threads,
	}
}
