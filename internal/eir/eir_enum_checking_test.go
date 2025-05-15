package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr2EirCheckingEnabled(t *testing.T) {
	checking := Str2EirChecking("enabled")

	assert.Equal(t, EIREnabled, checking)
}

func TestStr2EirCheckingDisabled(t *testing.T) {
	checking := Str2EirChecking("disabled")

	assert.Equal(t, EIRDisabled, checking)
}

func TestStr2EirCheckingMandatory(t *testing.T) {
	checking := Str2EirChecking("mandatory")

	assert.Equal(t, EIRMandatory, checking)
}

func TestStr2EirCheckinDisabledDefault(t *testing.T) {
	checking := Str2EirChecking("")

	assert.Equal(t, EIRDisabled, checking)
}

func TestStr2EirCheckinUnknown(t *testing.T) {
	checking := Str2EirChecking("hello")

	assert.Equal(t, EIRUnknown, checking)
}

func TestEirChecking2StrMandatory(t *testing.T) {
	mandatory := EirChecking2Str(EIRMandatory)

	assert.Equal(t, "mandatory", mandatory)
}
