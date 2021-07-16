package debug

import (
	"fmt"

	"github.com/spf13/viper"
)

//Println debug输出
func Println(a ...interface{}) (n int, err error) {
	if viper.GetBool("debug") == true {
		fmt.Printf("[DEBUG]")
		return fmt.Println(a...)
	}
	return 0, nil
}
