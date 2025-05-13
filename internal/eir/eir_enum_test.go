package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr2EirEnabled(t *testing.T) {
	enabled := Str2EirChecking("enabled")

	assert.Equal(t, EIREnabled, enabled)
}

func TestStr2EirDisabled(t *testing.T) {
	disabled := Str2EirChecking("disabled")

	assert.Equal(t, EIRDisabled, disabled)
}

func TestStr2EirMandatory(t *testing.T) {
	mandatory := Str2EirChecking("mandatory")

	assert.Equal(t, EIRMandatory, mandatory)
}

func TestEir2StrMandatory(t *testing.T) {
	mandatory := EirChecking2Str(EIRMandatory)

	assert.Equal(t, "mandatory", mandatory)
}
