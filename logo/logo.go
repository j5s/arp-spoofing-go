package logo

import (
	"bufio"
	"fmt"
	"os"
)

//LogoFile 开始时显示的Logo
var LogoFile string = "./logo/logo.txt"

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
