package version

import (
	"fmt"
	"os"
)

var (
	// BuildVersion : 版本号
	BuildVersion string = "1.0.1"
	// BuildTime : 版本构建时间
	BuildTime string = "2019.3.24"
	// BuildName : 构建工程名称
	BuildName string = "Magic Game"
	// Build 构建的版本 debug | release
	Build string
	// CommitID : 程序构建校验码
	CommitID string = "9aada4486c154195831f2e84c3036caf"
)

// Show : 显示版本信息
func Show() {
	fmt.Printf("Build:\t%s\n", Build)
	fmt.Printf("Build name:\t%s\n", BuildName)
	fmt.Printf("Build ver:\t%s\n", BuildVersion)
	fmt.Printf("Build time:\t%s\n", BuildTime)
	fmt.Printf("Program Commit ID:\t%s\n", CommitID)
	os.Exit(0)
}
