package version

import (
	"fmt"
	"os"
)

var (
	// BuildVersion : 版本号
	BuildVersion string
	// BuildTime : 版本构建时间
	BuildTime string
	// BuildName : 构建工程名称
	BuildName string
	// Build 构建的版本 debug | release
	Build string
	// CommitID : 程序构建校验码
	CommitID string
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
