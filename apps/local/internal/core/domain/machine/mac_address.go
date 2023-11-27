package machine

import (
	"errors"
	"regexp"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
)

const MacAddressPattern = `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`

var (
	ErrInvalidMacAddress = errors.New("invalid mac address")
)

type MacAddress string

func MacAddressFromString(v string) (MacAddress, error) {
	macPattern := regexp.MustCompile(MacAddressPattern)
	if !macPattern.MatchString(v) {
		return "", camperror.New(camperror.ErrValidation, ErrInvalidMacAddress)
	}

	return MacAddress(v), nil
}

func (m MacAddress) String() string {
	return string(m)
}
