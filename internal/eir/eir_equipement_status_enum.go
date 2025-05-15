package context

type EirEquipementStatus int

const (
	EIRWhitelisted   EirEquipementStatus = 0
	EIRGreylisted    EirEquipementStatus = 1
	EIRBlacklisted   EirEquipementStatus = 2
	EIRUnknownlisted EirEquipementStatus = 3
)

func Str2EirEquipementStatus(eirEquipementStatus string) EirEquipementStatus {
	switch eirEquipementStatus {
	case "WHITELISTED":
		return EIRWhitelisted
	case "GREYLISTED":
		return EIRGreylisted
	case "BLACKLISTED":
		return EIRBlacklisted
	default:
		return EIRUnknownlisted
	}
}

func EirEquipementStatus2Str(eirEquipementStatus EirEquipementStatus) string {
	switch eirEquipementStatus {
	case EIRWhitelisted:
		return "WHITELISTED"
	case EIRGreylisted:
		return "GREYLISTED"
	case EIRBlacklisted:
		return "BLACKLISTED"
	default:
		return ""
	}
}
