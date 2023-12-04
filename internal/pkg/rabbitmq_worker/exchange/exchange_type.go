package exchange

import "fmt"

type ExchangeType int

const (
	ExchangeTypeDirect ExchangeType = iota
	ExchangeTypeFanout
	ExchangeTypeTopic
	ExchangeTypeHeaders
)

func (e ExchangeType) String() string {
	return [...]string{"direct", "fanout", "topic", "headers"}[e]
}

func ExchangeTypeFromString(s string) (ExchangeType, error) {
	switch s {
	case "direct":
		return ExchangeTypeDirect, nil
	case "fanout":
		return ExchangeTypeFanout, nil
	case "topic":
		return ExchangeTypeTopic, nil
	case "headers":
		return ExchangeTypeHeaders, nil
	default:
		return ExchangeTypeDirect, fmt.Errorf("invalid exchange type %s", s)
	}
}
