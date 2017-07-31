// +build windows

package net

import (
	"strings"
	"time"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/StackExchange/wmi"
)

const (
	netQuery  = "SELECT Index,DHCPEnabled,DHCPLeaseExpires,DHCPLeaseObtained,DNSServerSearchOrder,IPEnabled,IPAddress,IPSubnet,DefaultIPGateway,DHCPServer,MACAddress,WINSPrimaryServer,WINSSecondaryServer,Description FROM Win32_NetworkAdapterConfiguration"
	netQuery2 = "SELECT Name,Manufacturer,Index FROM Win32_NetworkAdapter"
)

// win32NetworkAdapterConfiguration struct represents a network adapter configuratons
type win32NetworkAdapterConfiguration struct {
	Index                uint32
	DHCPEnabled          bool
	IPEnabled            bool
	DHCPLeaseExpires     *time.Time
	DHCPLeaseObtained    *time.Time
	DNSServerSearchOrder *[]string
	IPAddress            *[]string
	IPSubnet             *[]string
	DefaultIPGateway     *[]string
	DHCPServer           *string
	MACAddress           *string
	WINSPrimaryServer    *string
	WINSSecondaryServer  *string
	Description          string
}

// win32NetworkAdapter struct represents a network adapter
type win32NetworkAdapter struct {
	Name         string
	Manufacturer string
	Index        uint32
}

// Info returns network information for Windows using WMI
func Info() ([]asset.AssetNetwork, error) {
	var dst []win32NetworkAdapterConfiguration
	err := wmi.Query(netQuery, &dst)
	if err != nil {
		return nil, err
	}
	var dst2 []win32NetworkAdapter
	wmi.Query(netQuery2, &dst2)

	netArray := getAssetNetwork(dst, dst2)
	return netArray, nil
}

func getAssetNetwork(dst []win32NetworkAdapterConfiguration, dst2 []win32NetworkAdapter) []asset.AssetNetwork {
	netArray := make([]asset.AssetNetwork, len(dst))
	for i, v := range dst {
		var manuf string
		name := v.Description
		for _, w := range dst2 {
			if v.Index == w.Index {
				manuf = w.Manufacturer
				name = w.Name
				break
			}
		}
		var ipv4, ipv6, subnet, gateway, dhcpsvr = "0.0.0.0", "::", "0.0.0.0", "0.0.0.0", "0.0.0.0"
		var ipv4s, ipv6s []string

		getIPAddress(v.IPAddress, &ipv4s, &ipv6s, &ipv4, &ipv6)

		var gateways []string
		getArrayValue(v.DefaultIPGateway, &gateways)
		if len(gateways) > 0 {
			gateway = gateways[0]
		}
		var subnets []string
		getArrayValue(v.IPSubnet, &subnets)
		if len(subnets) > 0 {
			subnet = subnets[0]
		}
		var dnsservers []string
		getArrayValue(v.DNSServerSearchOrder, &dnsservers)

		var mac string
		getStringValue(v.MACAddress, &mac)
		getStringValue(v.DHCPServer, &dhcpsvr)

		var lobt, lexp time.Time
		getDateValue(v.DHCPLeaseObtained, &lobt)
		getDateValue(v.DHCPLeaseExpires, &lexp)

		var winsp, winss = "0.0.0.0", "0.0.0.0"
		getStringValue(v.WINSPrimaryServer, &winsp)
		getStringValue(v.WINSSecondaryServer, &winss)

		adapter := asset.AssetNetwork{
			Vendor:              manuf,
			Product:             name,
			LogicalName:         name,
			DhcpEnabled:         v.DHCPEnabled,
			DhcpServer:          dhcpsvr,
			DhcpLeaseObtained:   lobt,
			DhcpLeaseExpires:    lexp,
			DnsServers:          dnsservers,
			IPEnabled:           v.IPEnabled,
			IPv4:                ipv4,
			IPv4List:            ipv4s,
			IPv6:                ipv6,
			IPv6List:            ipv6s,
			SubnetMask:          subnet,
			SubnetMasks:         subnets,
			DefaultIPGateway:    gateway,
			DefaultIPGateways:   gateways,
			MacAddress:          mac,
			WinsPrimaryServer:   winsp,
			WinsSecondaryServer: winss,
		}
		netArray[i] = adapter
	}

	return netArray
}

func getIPAddress(ptrIPAdd, ipv4s, ipv6s *[]string, ipv4, ipv6 *string) {
	if ptrIPAdd != nil {
		for _, value := range *ptrIPAdd {
			if strings.Contains(value, ":") {
				*ipv6s = append(*ipv6s, value)
			} else {
				*ipv4s = append(*ipv4s, value)
			}
		}
		if len(*ipv4s) > 0 {
			*ipv4 = (*ipv4s)[0]
		}
		if len(*ipv6s) > 0 {
			*ipv6 = (*ipv6s)[0]
		}
	}
}
func getArrayValue(ptr, str *[]string) {
	if ptr != nil {
		*str = *ptr
	}
}
func getStringValue(ptr, str *string) {
	if ptr != nil && len(*ptr) > 0 {
		*str = *ptr
	}
}
func getDateValue(ptr, str *time.Time) {
	if ptr != nil {
		*str = *ptr
	}
}
