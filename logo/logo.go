package logo

import (
	"bufio"
	"fmt"
	"os"
)

var HOME string = os.Getenv("HOME")

//LogoFile 开始时显示的Logo
var LogoFile string = fmt.Sprintf("%s/src/arp-spoofing/logo/logo.txt", HOME)

//Show 显示Logo
func Show(logofile string) {
	file, err := os.Open(logofile)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
