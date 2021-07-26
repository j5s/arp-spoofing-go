package privileges

import (
	"os"
	"runtime"
)

//Check 检查用户权限
func Check() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	return os.Getuid() != 0
}
