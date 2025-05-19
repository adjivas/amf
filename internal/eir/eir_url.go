package context

import (
	"errors"
	"math/rand"
	"net/netip"

	"github.com/free5gc/openapi/models"
)

func getUriFromIpStr(scheme models.UriScheme, addrStr string, port int32) (string, error) {
	addr, err := netip.ParseAddr(addrStr)
	if err != nil {
		return "", err
	}

	portSelected := uint16(port)
	if port == 0 {
		switch scheme {
		case models.UriScheme_HTTPS:
			portSelected = 443
		case models.UriScheme_HTTP:
			portSelected = 80
		default:
			return "", errors.New("no port for found")
		}
	}

	bindAddr := netip.AddrPortFrom(addr, portSelected).String()

	eirUrl := string(scheme) + "://" + bindAddr

	return eirUrl, nil
}

func GetServiceNfUriFromIP(nfProfile *models.NrfNfManagementNfProfile, service models.NrfNfManagementNfService) (string, error) {
	// Get IP from NFService
	ipEndPointsLen := len(service.IpEndPoints)
	ipEndPointSelect := rand.Intn(ipEndPointsLen)
	ipEndPoint := service.IpEndPoints[ipEndPointSelect]

	if ipEndPoint.Ipv6Address != "" {
		return getUriFromIpStr(service.Scheme, ipEndPoint.Ipv6Address, ipEndPoint.Port)
	}
	if ipEndPoint.Ipv4Address != "" {
		return getUriFromIpStr(service.Scheme, ipEndPoint.Ipv4Address, ipEndPoint.Port)
	}
	// Get IP from NFProfile parent's NFService
	ipAddr6Len := len(nfProfile.Ipv4Addresses)
	ipAddr6Select := rand.Intn(ipAddr6Len)
	ipAddr6 := nfProfile.Ipv4Addresses[ipAddr6Select]
	if ipAddr6 != "" {
		return getUriFromIpStr(service.Scheme, ipAddr6, ipEndPoint.Port)
	}

	ipAddr4Len := len(nfProfile.Ipv4Addresses)
	ipAddr4Select := rand.Intn(ipAddr4Len)
	ipAddr4 := nfProfile.Ipv4Addresses[ipAddr4Select]
	if ipAddr4 != "" {
		return getUriFromIpStr(service.Scheme, ipAddr4, ipEndPoint.Port)
	}

	return "", errors.New("no uri found")
}

func GetServiceNfUri(nfProfile *models.NrfNfManagementNfProfile) (string, error) {
	var nfUri string
	for index := range nfProfile.NfServices {
		service := nfProfile.NfServices[index]
		if service.Fqdn != "" {
			nfUri = string(service.Scheme) + "://" + service.Fqdn
		} else if nfProfile.Fqdn != "" {
			nfUri = string(service.Scheme) + "://" + nfProfile.Fqdn
		} else if len(service.IpEndPoints) != 0 {
			nfUri, _ = GetServiceNfUriFromIP(nfProfile, service)
		}
		if nfUri != "" {
			break
		}
	}
	if nfUri == "" {
		return "", errors.New("no uri found")
	}
	return nfUri, nil
}
