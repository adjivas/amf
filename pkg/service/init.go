package service

import (
	"context"
	"errors"
	"io"
	"os"
	"runtime/debug"
	"sync"

	amf_context "github.com/free5gc/amf/internal/context"
	eir "github.com/free5gc/amf/internal/eir"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/ngap"
	ngap_message "github.com/free5gc/amf/internal/ngap/message"
	ngap_service "github.com/free5gc/amf/internal/ngap/service"
	"github.com/free5gc/amf/internal/sbi"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/amf/internal/sbi/processor"
	callback "github.com/free5gc/amf/internal/sbi/processor/notifier"
	"github.com/free5gc/amf/pkg/app"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/openapi/models"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
	"github.com/sirupsen/logrus"
)

type AmfAppInterface interface {
	app.App
	consumer.ConsumerAmf
	Consumer() *consumer.Consumer
	Processor() *processor.Processor
}

var AMF AmfAppInterface

type AmfApp struct {
	AmfAppInterface

	cfg               *factory.Config
	amfCtx            *amf_context.AMFContext
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	eirSubscriptionID string

	processor *processor.Processor
	consumer  *consumer.Consumer
	sbiServer *sbi.Server
}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*AmfApp, error) {
	amf := &AmfApp{
		cfg: cfg,
	}
	amf.SetLogEnable(cfg.GetLogEnable())
	amf.SetLogLevel(cfg.GetLogLevel())
	amf.SetReportCaller(cfg.GetLogReportCaller())

	consumer, err := consumer.NewConsumer(amf)
	if err != nil {
		return amf, err
	}
	amf.consumer = consumer

	processor, err_p := processor.NewProcessor(amf)
	if err_p != nil {
		return amf, err_p
	}
	amf.processor = processor

	amf.ctx, amf.cancel = context.WithCancel(ctx)
	amf.amfCtx = amf_context.GetSelf()

	if amf.sbiServer, err = sbi.NewServer(amf, tlsKeyLogPath); err != nil {
		return nil, err
	}

	AMF = amf

	return amf, nil
}

func (a *AmfApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}

	a.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (a *AmfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}

	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}

	a.cfg.SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *AmfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}

	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *AmfApp) Start() {
	self := a.Context()
	amf_context.InitAmfContext(self)

	ngapHandler := ngap_service.NGAPHandler{
		HandleMessage:         ngap.Dispatch,
		HandleNotification:    ngap.HandleSCTPNotification,
		HandleConnectionError: ngap.HandleSCTPConnError,
	}

	sctpConfig := ngap_service.NewSctpConfig(factory.AmfConfig.GetSctpConfig())
	ngap_service.Run(a.Context().NgapIpList, a.Context().NgapPort, ngapHandler, sctpConfig)
	logger.InitLog.Infoln("Server started")

	a.wg.Add(1)
	go a.listenShutdownEvent()

	var profile models.NrfNfManagementNfProfile
	if profileTmp, err := a.Consumer().BuildNFInstance(self); err != nil {
		logger.InitLog.Error("Build AMF Profile Error")
	} else {
		profile = profileTmp
	}

	_, nfId, err_reg := a.Consumer().SendRegisterNFInstance(a.ctx, a.Context().NrfUri, a.Context().NfId, &profile)
	if err_reg != nil {
		logger.InitLog.Warnf("Send Register NF Instance failed: %+v", err_reg)
	} else {
		a.Context().NfId = nfId
	}

	// Init Eir
	if a.Context().EIRChecking == eir.EIREnabled || a.Context().EIRChecking == eir.EIRMandatory {
		EIRRegistrationInfo, err := a.SearchEirInstance()
		if err != nil {
			logger.MainLog.Warnf("Search Eir instance failed %+v", err)
		} else {
			a.Context().EIRRegistrationInfo = EIRRegistrationInfo
			logger.InitLog.Infof("Select the Eir instance [%+v] from [%+v]", EIRRegistrationInfo.EIRApiPrefix, EIRRegistrationInfo.NfInstanceUri)
		}

		uriAmf := a.Context().GetIPUri()
		logger.InitLog.Infof("Binding addr: [%+v]", uriAmf)

		a.createEirSubscriptionProcedure(EIRRegistrationInfo.NfInstanceUri, uriAmf)
	}

	if err := a.sbiServer.Run(context.Background(), &a.wg); err != nil {
		logger.MainLog.Fatalf("Run SBI server failed: %+v", err)
	}
	a.WaitRoutineStopped()
}

func (a *AmfApp) SearchEirInstance() (amf_context.EIRRegistrationInfo, error) {
	NrfUri := a.Context().NrfUri
	param := Nnrf_NFDiscovery.SearchNFInstancesRequest{}
	resp, err := a.consumer.SendSearchNFInstances(NrfUri, models.NrfNfManagementNfType__5_G_EIR, models.NrfNfManagementNfType_AMF, &param)

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
	nfInstance := resp.NfInstances[0]

	if len(nfInstance.NfServices) <= 0 {
		return amf_context.EIRRegistrationInfo{
			NfInstanceUri: "",
			EIRApiPrefix:  "",
		}, errors.New("Not any NfServices were found")
	}
	nfServices := nfInstance.NfServices[0]

	prefix, errPrefix := eir.PrefixFromNfDiscoveryProfile(nfServices)
	if errPrefix != nil {
		logger.EIRLog.Warnf("The EIR notification is ignored because it's NfProfile is incorrect [%+v]", errPrefix)
	}
	nrfUri := factory.AmfConfig.GetNrfUri()
	return amf_context.EIRRegistrationInfo{
		NfInstanceUri: nrfUri + "/nnrf-nfm/v1/nf-instances/" + nfInstance.NfInstanceId,
		EIRApiPrefix:  prefix,
	}, nil
}

func (a *AmfApp) createEirSubscriptionProcedure(NfInstanceIdEir string, uriAmf string) {
	subscriptionData := Nnrf_NFManagement.CreateSubscriptionRequest{
		NrfNfManagementSubscriptionData: &models.NrfNfManagementSubscriptionData{
			NfStatusNotificationUri: uriAmf + "/namf-callback/v1/nnrf-nfm/v1",
			SubscrCond: &models.SubscrCond{
				NfType:       string(models.NrfNfManagementNfType__5_G_EIR),
				ServiceName:  models.ServiceName_N5G_EIR_EIC,
				NfInstanceId: NfInstanceIdEir,
			},
		},
	}
	uri := a.Context().NrfUri
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		logger.MainLog.Errorf("Failed to get NRF token %+v", err)
	}

	response, err := client.SubscriptionsCollectionApi.CreateSubscription(ctx, &subscriptionData)
	if err != nil {
		logger.MainLog.Errorf("Send Subscriptions nRF Eir failed %+v", err)
	} else {
		logger.InitLog.Infof("Registered Subscriptions nRF Eir %+v", response.NrfNfManagementSubscriptionData.SubscriptionId)
		a.eirSubscriptionID = response.NrfNfManagementSubscriptionData.SubscriptionId
	}
}

func (a *AmfApp) removeEirSubscriptionProcedure() {
	if eirSubscriptionID := a.eirSubscriptionID; eirSubscriptionID != "" {
		uri := a.Context().NrfUri
		configuration := Nnrf_NFManagement.NewConfiguration()
		configuration.SetBasePath(uri)
		client := Nnrf_NFManagement.NewAPIClient(configuration)

		request := Nnrf_NFManagement.RemoveSubscriptionRequest{
			SubscriptionID: &eirSubscriptionID,
		}

		ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
		if err != nil {
			logger.MainLog.Errorf("Failed to get NRF token %+v", err)
		}
		response, err := client.SubscriptionIDDocumentApi.RemoveSubscription(ctx, &request)
		if err != nil {
			logger.MainLog.Errorf("Send RemoveSubscription nRF Eir failed %+v", err)
		} else {
			logger.InitLog.Infof("RemoveSubscription nRF Eir %+v", response)
		}
	}
}

// Used in AMF planned removal procedure
func (a *AmfApp) Terminate() {
	a.cancel()
}

func (a *AmfApp) Config() *factory.Config {
	return a.cfg
}

func (a *AmfApp) Context() *amf_context.AMFContext {
	return a.amfCtx
}

func (a *AmfApp) CancelContext() context.Context {
	return a.ctx
}

func (a *AmfApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *AmfApp) Processor() *processor.Processor {
	return a.processor
}

func (a *AmfApp) listenShutdownEvent() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		a.wg.Done()
	}()

	<-a.ctx.Done()
	a.terminateProcedure()
}

func (a *AmfApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Stop()
	}
}

func (a *AmfApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("AMF App is terminated")
}

func (a *AmfApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating AMF...")
	a.CallServerStop()

	// deregister with NRF
	a.removeEirSubscriptionProcedure()
	problemDetails, err_deg := a.Consumer().SendDeregisterNFInstance()
	if problemDetails != nil {
		logger.MainLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
	} else if err_deg != nil {
		logger.MainLog.Errorf("Deregister NF instance Error[%+v]", err_deg)
	} else {
		logger.MainLog.Infof("[AMF] Deregister from NRF successfully")
	}

	// TODO: forward registered UE contexts to target AMF in the same AMF set if there is one

	// ngap
	// send AMF status indication to ran to notify ran that this AMF will be unavailable
	logger.MainLog.Infof("Send AMF Status Indication to Notify RANs due to AMF terminating")
	amfSelf := a.Context()
	unavailableGuamiList := ngap_message.BuildUnavailableGUAMIList(amfSelf.ServedGuamiList)
	amfSelf.AmfRanPool.Range(func(key, value interface{}) bool {
		ran := value.(*amf_context.AmfRan)
		ngap_message.SendAMFStatusIndication(ran, unavailableGuamiList)
		return true
	})
	ngap_service.Stop()
	callback.SendAmfStatusChangeNotify((string)(models.StatusChange_UNAVAILABLE), amfSelf.ServedGuamiList)
}
