package context

import (
	"errors"
	"net/netip"

	"github.com/free5gc/openapi/models"
)

func PrefixFromNfDiscoveryProfile(nfService models.NrfNfDiscoveryNfService) (string, error) {
	if len(nfService.IpEndPoints) <= 0 {
		return "", errors.New("The nfServices doesn't have a IpEndPoint ")
	}
	ipEndPoint := nfService.IpEndPoints[0]

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

	eirUrl := string(nfService.Scheme) + "://" + bindAddr
	return eirUrl, nil
}

func PrefixFromNfProfile(nfService models.NrfNfManagementNfService) (string, error) {
	if len(nfService.IpEndPoints) <= 0 {
		return "", errors.New("The nfServices doesn't have a IpEndPoint ")
	}

	ipEndPoint := nfService.IpEndPoints[0]

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

	eirUrl := string(nfService.Scheme) + "://" + bindAddr + "/nnrf-nfm/v1/nf-instances/"

	return eirUrl, nil
}
