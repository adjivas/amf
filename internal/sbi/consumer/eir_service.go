package consumer

import (
	"sync"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/openapi"
	Neir_EIRSelection "github.com/free5gc/openapi/eir/EIRService"
	"github.com/free5gc/openapi/models"
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

func (s *neirService) GetEquipementStatus(uri string, imei string) (*Neir_EIRSelection.EIREquipementStatusGetResponse, error) {
	client := s.getEIRSelectionClient(uri)
	if client == nil {
		return nil, openapi.ReportError("eir not found")
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_N5G_EIR_EIC, models.NrfNfManagementNfType__5_G_EIR)
	if err != nil {
		return nil, err
	}

	return client.EIREquipementStatusApi.EIREquipementStatusGet(ctx, imei)
}
