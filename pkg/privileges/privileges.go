package privileges

import "os"

//Check 检查用户权限
func Check() bool {
	return os.Getuid() != 0
}
