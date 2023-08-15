package test

import (
	"fmt"
	"strings"
	"testing"
)

// go test -timeout 30s -run ^TestOk github.com/hootrhino/rulex/test -v -count=1
/*
*
* 读写ubuntu18的配置
*
 */
type NetworkInterface struct {
	Name        string
	Interface   string
	Address     string
	Netmask     string
	Gateway     string
	DNS         []string
	DHCPEnabled bool
}

func writeInterfaceConfig(filePath string, ifaces []NetworkInterface) {
	configText := ""
	for _, iface := range ifaces {
		configLines := []string{
			"auto lo",
			"iface lo inet loopback",
			fmt.Sprintf("auto %s", iface.Name),
			fmt.Sprintf("iface %s inet %s", iface.Interface, getInetType(iface.DHCPEnabled)),
		}

		if !iface.DHCPEnabled {
			configLines = append(configLines, fmt.Sprintf("    address %s", iface.Address))
			configLines = append(configLines, fmt.Sprintf("    netmask %s", iface.Netmask))
			configLines = append(configLines, fmt.Sprintf("    gateway %s", iface.Gateway))
			configLines = append(configLines, fmt.Sprintf("    dns-nameservers %s", strings.Join(iface.DNS, " ")))
		}
		configText += strings.Join(configLines, "\n")
	}
	fmt.Println(configText)

	// return os.WriteFile(filePath, []byte(configText), 0644)
}

func getInetType(dhcpEnabled bool) string {
	if dhcpEnabled {
		return "dhcp"
	}
	return "static"
}
func TestReadNetplanConfig(t *testing.T) {

}
func TestWriteNetplanConfig(t *testing.T) {

}

/*
*
* 读写ubuntu16的配置
*
 */
func TestReadEtcNetConfig(t *testing.T) {

}

// go test -timeout 30s -run ^TestWriteFEtcNetConfig github.com/hootrhino/rulex/test -v -count=1

func TestWriteFEtcNetConfig(t *testing.T) {
	interfaceConfig := NetworkInterface{
		Name:        "eth0",
		Interface:   "eth0",
		Address:     "192.168.1.100",
		Netmask:     "255.255.255.0",
		Gateway:     "192.168.1.1",
		DNS:         []string{"8.8.8.8", "114.114.114.114"},
		DHCPEnabled: false,
	}

	writeInterfaceConfig("/etc/network/interfaces", []NetworkInterface{interfaceConfig, interfaceConfig})
	fmt.Println("Configuration written successfully.")
}
