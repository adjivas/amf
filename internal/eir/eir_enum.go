package context

import (
	"fmt"

	"github.com/free5gc/aper"
)

const (
	EIREnabled   aper.Enumerated = 0
	EIRDisabled  aper.Enumerated = 1
	EIRMandatory aper.Enumerated = 2
)

type EirChecking struct {
	Value aper.Enumerated `aper:"enabled:0,disabled:1,mandatory:2`
}

func Str2EirChecking(eirChecking string) aper.Enumerated {
	switch eirChecking {
	case "enabled":
		return EIREnabled
	case "disabled":
		return EIRDisabled
	case "mandatory":
		return EIRMandatory
	case "":
		return EIRDisabled
	default:
		panic(fmt.Sprintf("Unknown EirChecking %s", eirChecking))
	}
}

func EirChecking2Str(eirChecking aper.Enumerated) string {
	switch eirChecking {
	case EIREnabled:
		return "enabled"
	case EIRDisabled:
		return "disabled"
	case EIRMandatory:
		return "mandatory"
	default:
		panic(fmt.Sprintf("Unknown EirChecking %d", eirChecking))
	}
}
