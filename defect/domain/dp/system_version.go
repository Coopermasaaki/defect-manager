package dp

import (
	"errors"
)

var MaintainVersion = make(map[SystemVersion]bool)

type systemVersion string

type SystemVersion interface {
	String() string
}

func Init(maintainVersion []string) {
	for _, version := range maintainVersion {
		MaintainVersion[systemVersion(version)] = true
	}
}

func NewSystemVersion(s string) (SystemVersion, error) {
	// MaintainVersion is not used for validation because
	// there is an error reading old data from the database when maintainVersion changes
	if s == "" {
		return nil, errors.New("invalid system version")
	}

	return systemVersion(s), nil
}

func (s systemVersion) String() string {
	return string(s)
}
