package consumer

import (
	"sync"
	"errors"

	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/amf/pkg/factory"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Neir_EIRSelection "github.com/free5gc/openapi/eir/EIRService"
	amf_context "github.com/free5gc/amf/internal/context"
)

type neirService struct {
	consumer *Consumer

	EIRSelectionMu sync.RWMutex

	EIRSelectionClients map[string]*Neir_EIRSelection.APIClient
}

func (s *neirService) getEIRSelectionClient(uri string) *Neir_EIRSelection.APIClient {
	if uri == "" {
		return nil
	}
	s.EIRSelectionMu.RLock()
	client, ok := s.EIRSelectionClients[uri]
	if ok {
		s.EIRSelectionMu.RUnlock()
		return client
	}

	configuration := Neir_EIRSelection.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Neir_EIRSelection.NewAPIClient(configuration)

	s.EIRSelectionMu.RUnlock()
	s.EIRSelectionMu.Lock()
	defer s.EIRSelectionMu.Unlock()
	s.EIRSelectionClients[uri] = client
	return client
}

func (s *neirService) GetEquipmentStatus(uri string, imei string) (*Neir_EIRSelection.EIREquipmentStatusGetResponse, error) {
	client := s.getEIRSelectionClient(uri)
	if client == nil {
		return nil, openapi.ReportError("eir not found")
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_N5G_EIR_EIC, models.NrfNfManagementNfType__5_G_EIR)
	if err != nil {
		return nil, err
	}

	return client.EIREquipmentStatusApi.EIREquipmentStatusGet(ctx, imei)
}

func SearchEirInstance(consumer *Consumer) (amf_context.EIRRegistrationInfo, error) {
	NrfUri := amf_context.GetSelf().NrfUri
	param := Nnrf_NFDiscovery.SearchNFInstancesRequest{}
	resp, err := consumer.SendSearchNFInstances(NrfUri, models.NrfNfManagementNfType__5_G_EIR, models.NrfNfManagementNfType_AMF, &param)

	if err != nil {
		logger.MainLog.Errorf("Send Search NF Instances 5_G_EIR failed: %+v", err)
		return amf_context.EIRRegistrationInfo{
			NfInstanceUri: "",
			EIRApiPrefix:  "",
		}, err
	}

	if len(resp.NfInstances) <= 0 {
		return amf_context.EIRRegistrationInfo{
			NfInstanceUri: "",
			EIRApiPrefix:  "",
		}, errors.New("Not any NfInstances were found")
	}

	nfProfile, eirUri, errProfile := openapi.GetServiceNfProfileAndUri(resp.NfInstances, models.ServiceName_N5G_EIR_EIC)
	if errProfile != nil {
		logger.EIRLog.Warnf("The EIR notification is ignored because it's NfProfile is incorrect [%+v]", errProfile)
	}
	nrfUri := factory.AmfConfig.GetNrfUri()
	return amf_context.EIRRegistrationInfo{
		NfInstanceUri: nrfUri + "/nnrf-nfm/v1/nf-instances/" + nfProfile.NfInstanceId,
		EIRApiPrefix:  eirUri,
	}, nil
}
