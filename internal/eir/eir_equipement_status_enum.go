package context

type EirEquipmentStatus int

const (
	EIRWhitelisted   EirEquipmentStatus = 0
	EIRGreylisted    EirEquipmentStatus = 1
	EIRBlacklisted   EirEquipmentStatus = 2
	EIRUnknownlisted EirEquipmentStatus = 3
)

func Str2EirEquipmentStatus(eirEquipmentStatus string) EirEquipmentStatus {
	switch eirEquipmentStatus {
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

func EirEquipmentStatus2Str(eirEquipmentStatus EirEquipmentStatus) string {
	switch eirEquipmentStatus {
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
