package machine

type Cpu struct {
	Cores        int
	Architecture Architecture
	Model        *string
	Vendor       *string
}

func NewCpu(cores int, architecture Architecture, model *string, vendor *string) Cpu {
	return Cpu{
		Cores:        cores,
		Architecture: architecture,
		Model:        model,
		Vendor:       vendor,
	}
}
