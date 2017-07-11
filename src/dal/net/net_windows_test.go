// +build windows

package net

import (
	"testing"
	"time"
)

func TestGetAssetNetwork(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now()
	dns := []string{"1.1.1.1", "2.2.2.2"}
	subnet := []string{"10.10.10.10", "20.20.20.20"}
	gateway := []string{"110.110.110.110", "220.220.220.220"}
	ipadd := []string{"169.254.60.58", "169.254.60.59", "fe80::b5f0:e4d0:3bde:3c3a"}
	winsp := "1.2.3.4"
	winss := "4.3.2.1"
	dhcp := "10.20.30.40"
	mac := "0A:00:27:00:00:15"

	a := []win32_NetworkAdapterConfiguration{
		win32_NetworkAdapterConfiguration{
			Description:          "Wireless Adapter",
			Index:                0,
			DHCPLeaseExpires:     &t1,
			DHCPLeaseObtained:    &t2,
			DNSServerSearchOrder: &dns,
			IPSubnet:             &subnet,
			WINSPrimaryServer:    &winsp,
			WINSSecondaryServer:  &winss,
			DHCPServer:           &dhcp,
			DefaultIPGateway:     &gateway,
			MACAddress:           &mac,
			IPAddress:            &ipadd,
		},
		win32_NetworkAdapterConfiguration{
			Description:          "Ethernet Adapter",
			Index:                1,
			DHCPLeaseExpires:     nil,
			DHCPLeaseObtained:    nil,
			DNSServerSearchOrder: nil,
			IPSubnet:             nil,
			WINSPrimaryServer:    nil,
			WINSSecondaryServer:  nil,
			DHCPServer:           nil,
			DefaultIPGateway:     nil,
			MACAddress:           nil,
			IPAddress:            nil,
		},
	}
	b := []win32_NetworkAdapter{
		win32_NetworkAdapter{
			Manufacturer: "Intel Corp",
			Index:        0,
		},
		win32_NetworkAdapter{
			Manufacturer: "Microsoft Corp",
			Index:        1,
		},
	}
	res := getAssetNetwork(a, b)
	if len(b) != len(a) {
		t.Errorf("Wrong setup")
	}
	if len(res) != len(a) {
		t.Errorf("Unexpected count %v, expected count %d ", len(res), len(a))
	}
	for i, v := range res {
		if v.Vendor != b[i].Manufacturer {
			t.Errorf("Wrong vendor %s, expected vendor %s ", v.Vendor, b[i].Manufacturer)
		}
		if v.IPEnabled != a[i].IPEnabled {
			t.Errorf("Wrong IPEnabled %v, expected IPEnabled %v ", v.IPEnabled, a[i].IPEnabled)
		}
		var testIP string
		if a[i].IPAddress != nil && len(*a[i].IPAddress) > 0 {
			testIP = (*a[i].IPAddress)[0]
		}

		if v.IPv4 != testIP {
			t.Errorf("Wrong IPv4 %v, expected IPv4 %v ", v.IPv4, testIP)
		}
	}
}
