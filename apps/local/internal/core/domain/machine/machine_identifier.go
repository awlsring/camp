package machine

import (
	"errors"
	"fmt"
	"regexp"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
)

const (
	MachineIdentifierPattern = "^m-"
	MachineIdentifierLength  = 36
)

var (
	ErrInvalidMachineIdentifier = errors.New("invalid machine identifier")
)

type Identifier string

func NewIdentifier(internalIdentifier InternalIdentifier) Identifier {
	return Identifier(fmt.Sprintf("m-%s", internalIdentifier.String()))
}

func IdentifierFromString(identifier string) (Identifier, error) {
	re := regexp.MustCompile(MachineIdentifierPattern)
	if !re.MatchString(identifier) {
		return "", ErrInvalidMachineIdentifier
	}

	if len(identifier) != MachineIdentifierLength {
		return "", camperror.New(camperror.ErrValidation, ErrInvalidMachineIdentifier)
	}

	return Identifier(identifier), nil
}

func (m Identifier) String() string {
	return string(m)
}
