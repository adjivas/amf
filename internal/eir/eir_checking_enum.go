package context

type EirChecking int

const (
	EIREnabled   EirChecking = 0
	EIRDisabled  EirChecking = 1
	EIRMandatory EirChecking = 2
	EIRUnknown   EirChecking = 3
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
		return EIRUnknown
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
		return ""
	}
}
