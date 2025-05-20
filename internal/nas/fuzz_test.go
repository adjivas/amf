//go:build go1.18
// +build go1.18

package nas_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	amf_nas "github.com/free5gc/amf/internal/nas"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/amf/pkg/service"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func FuzzHandleNAS(f *testing.F) {
	amfSelf := amf_context.GetSelf()
	amfSelf.ServedGuamiList = []models.Guami{{
		PlmnId: &models.PlmnIdNid{
			Mcc: "208",
			Mnc: "93",
		},
		AmfId: "cafe00",
	}}
	tai := models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "1",
	}
	amfSelf.SupportTaiLists = []models.Tai{tai}

	msg := nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)
	msg.RegistrationRequest = nasMessage.NewRegistrationRequest(nas.MsgTypeRegistrationRequest)
	reg := msg.RegistrationRequest
	reg.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	reg.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	reg.SetMessageType(nas.MsgTypeRegistrationRequest)
	reg.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	reg.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(7)
	reg.SetFOR(1)
	reg.SetRegistrationType5GS(nasMessage.RegistrationType5GSInitialRegistration)
	id := []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}
	reg.MobileIdentity5GS.SetLen(uint16(len(id)))
	reg.SetMobileIdentity5GSContents(id)
	reg.UESecurityCapability = nasType.NewUESecurityCapability(nasMessage.RegistrationRequestUESecurityCapabilityType)
	reg.UESecurityCapability.SetLen(2)
	reg.SetEA0_5G(1)
	reg.SetIA2_128_5G(1)
	buf, err := msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	deReg := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(
		nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	msg.DeregistrationRequestUEOriginatingDeregistration = deReg
	deReg.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deReg.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deReg.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	deReg.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	deReg.SetNasKeySetIdentifiler(7)
	deReg.SetSwitchOff(0)
	deReg.SetAccessType(nasMessage.AccessType3GPP)
	deReg.SetLen(uint16(len(id)))
	deReg.SetMobileIdentity5GSContents(id)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeServiceRequest)
	msg.ServiceRequest = nasMessage.NewServiceRequest(nas.MsgTypeServiceRequest)
	sr := msg.ServiceRequest
	sr.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	sr.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	sr.SetMessageType(nas.MsgTypeServiceRequest)
	sr.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	sr.SetNasKeySetIdentifiler(0)
	sr.SetServiceTypeValue(nasMessage.ServiceTypeSignalling)
	sr.TMSI5GS.SetLen(7)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	buf = append([]uint8{
		nasMessage.Epd5GSMobilityManagementMessage,
		nas.SecurityHeaderTypeIntegrityProtected,
		0, 0, 0, 0, 0,
	},
		buf...)
	f.Add(buf)

	f.Fuzz(func(t *testing.T, d []byte) {
		ue := new(amf_context.RanUe)
		ue.Ran = new(amf_context.AmfRan)
		ue.Ran.AnType = models.AccessType__3_GPP_ACCESS
		ue.Ran.Log = logger.NgapLog
		ue.Log = logger.NgapLog
		ue.Tai = tai
		ue.AmfUe = amfSelf.NewAmfUe("")
		amf_nas.HandleNAS(ue, ngapType.ProcedureCodeInitialUEMessage, d, true)
	})
}

func FuzzHandleNAS2(f *testing.F) {
	amfSelf := amf_context.GetSelf()
	amfSelf.ServedGuamiList = []models.Guami{{
		PlmnId: &models.PlmnIdNid{
			Mcc: "208",
			Mnc: "93",
		},
		AmfId: "cafe00",
	}}
	tai := models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "1",
	}
	amfSelf.SupportTaiLists = []models.Tai{tai}
	amfSelf.NrfUri = "test"

	msg := nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)
	msg.RegistrationRequest = nasMessage.NewRegistrationRequest(nas.MsgTypeRegistrationRequest)
	reg := msg.RegistrationRequest
	reg.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	reg.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	reg.SetMessageType(nas.MsgTypeRegistrationRequest)
	reg.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	reg.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(7)
	reg.SetFOR(1)
	reg.SetRegistrationType5GS(nasMessage.RegistrationType5GSInitialRegistration)
	id := []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}
	reg.MobileIdentity5GS.SetLen(uint16(len(id)))
	reg.SetMobileIdentity5GSContents(id)
	reg.UESecurityCapability = nasType.NewUESecurityCapability(nasMessage.RegistrationRequestUESecurityCapabilityType)
	reg.UESecurityCapability.SetLen(2)
	reg.SetEA0_5G(1)
	reg.SetIA2_128_5G(1)
	regPkt, err := msg.PlainNasEncode()
	require.NoError(f, err)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeIdentityResponse)
	msg.IdentityResponse = nasMessage.NewIdentityResponse(nas.MsgTypeIdentityResponse)
	ir := msg.IdentityResponse
	ir.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ir.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ir.SetMessageType(nas.MsgTypeIdentityResponse)
	ir.SetLen(uint16(len(id)))
	ir.SetMobileIdentityContents(id)
	buf, err := msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationResponse)
	msg.AuthenticationResponse = nasMessage.NewAuthenticationResponse(nas.MsgTypeAuthenticationResponse)
	ar := msg.AuthenticationResponse
	ar.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ar.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ar.SetMessageType(nas.MsgTypeAuthenticationResponse)
	ar.AuthenticationResponseParameter = nasType.NewAuthenticationResponseParameter(
		nasMessage.AuthenticationResponseAuthenticationResponseParameterType)
	ar.AuthenticationResponseParameter.SetLen(16)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationFailure)
	msg.AuthenticationFailure = nasMessage.NewAuthenticationFailure(nas.MsgTypeAuthenticationFailure)
	af := msg.AuthenticationFailure
	af.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	af.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	af.SetMessageType(nas.MsgTypeAuthenticationFailure)
	af.SetCauseValue(nasMessage.Cause5GMMSynchFailure)
	af.AuthenticationFailureParameter = nasType.NewAuthenticationFailureParameter(
		nasMessage.AuthenticationFailureAuthenticationFailureParameterType)
	af.SetLen(14)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeStatus5GMM)
	msg.Status5GMM = nasMessage.NewStatus5GMM(nas.MsgTypeStatus5GMM)
	st := msg.Status5GMM
	st.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	st.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	st.SetMessageType(nas.MsgTypeStatus5GMM)
	st.SetCauseValue(nasMessage.Cause5GMMProtocolErrorUnspecified)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	f.Fuzz(func(t *testing.T, d []byte) {
		ctrl := gomock.NewController(t)
		// m := app.NewMockApp(ctrl)
		m := service.NewMockAmfAppInterface(ctrl)
		c, errc := consumer.NewConsumer(m)
		service.AMF = m
		require.NoError(t, errc)
		m.EXPECT().
			Consumer().
			AnyTimes().
			Return(c)

		ue := new(amf_context.RanUe)
		ue.Ran = new(amf_context.AmfRan)
		ue.Ran.AnType = models.AccessType__3_GPP_ACCESS
		ue.Ran.Log = logger.NgapLog
		ue.Log = logger.NgapLog
		ue.Tai = tai
		ue.AmfUe = amfSelf.NewAmfUe("")
		amf_nas.HandleNAS(ue, ngapType.ProcedureCodeInitialUEMessage, regPkt, true)
		amfUe := ue.AmfUe
		amfUe.State[models.AccessType__3_GPP_ACCESS].Set(amf_context.Authentication)
		amfUe.RequestIdentityType = nasMessage.MobileIdentity5GSTypeSuci
		amfUe.AuthenticationCtx = &models.UeAuthenticationCtx{
			AuthType: models.AusfUeAuthenticationAuthType__5_G_AKA,
		}
		amf_nas.HandleNAS(ue, ngapType.ProcedureCodeUplinkNASTransport, d, false)
	})
}
