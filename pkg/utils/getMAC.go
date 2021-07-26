package utils

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

//extractDeviceID 从设备名中提取设备ID （为了适配windows）
func extractDeviceID(device string) (deviceID string) {
	var deviceNameRe = regexp.MustCompile(`^\\Device\\[^{}]+_{([^}]+)}$`)
	matches := deviceNameRe.FindStringSubmatch(device)
	if len(matches) < 1 {
		return ""
	}
	return matches[1]
}

//GetMAC 获取mac地址
func GetMAC(device string) (mac net.HardwareAddr, err error) {
	//windows
	if runtime.GOOS == "windows" {
		deviceID := extractDeviceID(device)
		Cmd := exec.Command("cmd", "/C", "getmac") //windows
		retCmd, err := Cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Cmd.CombinedOuput failed,err:", err)
			return mac, err
		}
		result := string(retCmd)
		lines := strings.Split(result, "\n")
		for _, line := range lines {
			parts := strings.Split(line, "   ")
			if len(parts) < 2 {
				continue
			}
			if deviceID == extractDeviceID(strings.TrimSpace(parts[1])) {
				return net.ParseMAC(parts[0])
			}
		}
		return mac, nil
	}
	//Linux,macOS
	iface, err := net.InterfaceByName(device)
	if err != nil {
		fmt.Println("net.InterfaceByName faield,err:", err)
		return mac, err
	}
	return iface.HardwareAddr, nil
}
