package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr2EirEquipmentStatusWhitelisted(t *testing.T) {
	status := Str2EirEquipmentStatus("WHITELISTED")

	assert.Equal(t, EIRWhitelisted, status)
}

func TestStr2EirEquipmentStatusGreylisted(t *testing.T) {
	status := Str2EirEquipmentStatus("GREYLISTED")

	assert.Equal(t, EIRGreylisted, status)
}

func TestStr2EirEquipmentStatusBlacklisted(t *testing.T) {
	status := Str2EirEquipmentStatus("BLACKLISTED")

	assert.Equal(t, EIRBlacklisted, status)
}

func TestEirEquipmentStatus2StrUnknownlisted(t *testing.T) {
	status := Str2EirEquipmentStatus("PINKLISTED")

	assert.Equal(t, EIRUnknownlisted, status)
}

func TestEirEquipmentStatus2StrWhitelisted(t *testing.T) {
	status := EirEquipmentStatus2Str(EIRWhitelisted)

	assert.Equal(t, "WHITELISTED", status)
}
