package context

import (
	"fmt"
)

type EirChecking int

const (
	EIREnabled   EirChecking = 0
	EIRDisabled  EirChecking = 1
	EIRMandatory EirChecking = 2
)

func Str2EirChecking(eirChecking string) EirChecking {
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

func EirChecking2Str(eirChecking EirChecking) string {
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
