package context

import (
	"errors"
	"net/netip"

	"github.com/free5gc/openapi/models"
)

func PrefixFromAnyProfile(ipEndPoints []models.IpEndPoint, scheme string) (string, error) {
	if len(ipEndPoints) <= 0 {
		return "", errors.New("The nfServices doesn't have a IpEndPoint ")
	}

	ipEndPoint := ipEndPoints[0]

	var addrStr string
	if Ipv6Address := ipEndPoint.Ipv6Address; Ipv6Address != "" {
		addrStr = Ipv6Address
	} else if Ipv4Address := ipEndPoint.Ipv4Address; Ipv4Address != "" {
		addrStr = Ipv4Address
	} else {
		return "", errors.New("The nfServices IpEndPoint doesn't have a IP address")
	}

	addr, err := netip.ParseAddr(addrStr)
	if err != nil {
		return "", err
	}

	bindAddr := netip.AddrPortFrom(addr, uint16(ipEndPoint.Port)).String()

	eirUrl := string(scheme) + "://" + bindAddr

	return eirUrl, nil
}

func PrefixFromNfDiscoveryProfile(nfService models.NrfNfDiscoveryNfService) (string, error) {
	return PrefixFromAnyProfile(nfService.IpEndPoints, string(nfService.Scheme))
}

func PrefixFromNfProfile(nfService models.NrfNfManagementNfService) (string, error) {
	return PrefixFromAnyProfile(nfService.IpEndPoints, string(nfService.Scheme))
}
