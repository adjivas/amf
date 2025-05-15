package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr2EirEquipementStatusWhitelisted(t *testing.T) {
	status := Str2EirEquipementStatus("WHITELISTED")

	assert.Equal(t, EIRWhitelisted, status)
}

func TestStr2EirEquipementStatusGreylisted(t *testing.T) {
	status := Str2EirEquipementStatus("GREYLISTED")

	assert.Equal(t, EIRGreylisted, status)
}

func TestStr2EirEquipementStatusBlacklisted(t *testing.T) {
	status := Str2EirEquipementStatus("BLACKLISTED")

	assert.Equal(t, EIRBlacklisted, status)
}

func TestEirEquipementStatus2StrUnknownlisted(t *testing.T) {
	status := Str2EirEquipementStatus("PINKLISTED")

	assert.Equal(t, EIRUnknownlisted, status)
}

func TestEirEquipementStatus2StrWhitelisted(t *testing.T) {
	status := EirEquipementStatus2Str(EIRWhitelisted)

	assert.Equal(t, "WHITELISTED", status)
}
