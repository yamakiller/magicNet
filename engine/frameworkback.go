package engine

/*import (
	"flag"
	"fmt"
	"io/ioutil"
	"magicNet/engine/logger"
	"magicNet/engine/monitor"
	"magicNet/engine/timer"
	"magicNet/engine/util"
	"magicNet/engine/hook"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	verMajor     = "1"
	verMinor     = "0"
	verPatch     = "0"
	verBuild     = "x64"
	verRevision  = "Beta"
	verCode      = "UTF8"
	verBuildTime = "2019-06-29 11:58"
)


// Framework : 主框架对象
type Framework struct {
	name       string
	configPath string
	loggerPath string
	loggerLv   string
}

var engineInitHook hook.InitializeHook

// SetEngineInitHook : 设置引擎初始/销毁Hook函数
func SetEngineInitHook(enHook hook.InitializeHook) {
  if engineInitHook == nil {
    engineInitHook = enHook
  }
}*/

// Start is Start system framework
/*func (fr *Framework) Start() int {

	isHelp := false
	isVers := false
  flag.BoolVar(&isHelp, "h", false, "out help informat")
	flag.BoolVar(&isHelp, "?", false, "out help informat")
	flag.BoolVar(&isVers, "v", false, "display version informat")
	flag.StringVar(&fr.configPath, "c", "./conf/magicnet.conf", "config full path")
	flag.StringVar(&fr.name, "n", "magic network base", "system name")
	flag.StringVar(&fr.loggerPath, "g", "", "system log file path")
	flag.StringVar(&fr.loggerLv, "l", "all", "system log level")
	flag.Parse()

	if isHelp {
		fr.usage()
		return 1
	} else if isVers {
		fr.version()
		return 1
	}

	monitor.Init()
	monitor.SetStateStart()

	logger.StartService(fr.loggerLv)

	logger.Redirect(fr.loggerPath)

	logger.Info(0, "loading env")

	if strings.Compare(fr.configPath, "") == 0 {
		logger.Error(0, "please set the configuration file, please enter -help=true to view the parameters.")
		return -1
	}

	if util.LoadEnv(fr.configPath) != 0 {
		return -1
	}

	logger.Info(0, "start %s ....", fr.name)
	if !monitor.StartService() {
		return -1
	}
	return fr.bootstrap()
}

// Loop framework mian loop
func (fr *Framework) Loop() {
	if !engineInitHook.Initialize() {
		return
	}

	monitor.SetStateRun()
	for !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
	}
	monitor.SetStateShutdown()


	timer.StopService()
	monitor.StopService()
	monitor.WaitSupper()
}

// Shutdown framework end
func (fr *Framework) Dispose() {
	engineInitHook.Finalize()
	logger.Info(0, "%s exit", fr.name)
	logger.SotpService()
	util.UnLoadEnv()
	monitor.SetStateIdle()
}

func (fr *Framework) bootstrap() int {
	timer.StartService()
	return 0
}

func (fr *Framework) version() {
	fmt.Printf("Version: %s.%s.%s %s %s\n", verMajor, verMinor, verPatch, verBuild, verRevision)
	fmt.Printf("Build Time: %s\n", verBuildTime)
	fmt.Printf("Go Version: %s %v\n", runtime.Version(), runtime.Compiler)
}

func (fr *Framework) usage() {
	f, err := os.Open("help.md")
	if err != nil {
		fmt.Printf("Error: open help.md file fail:%v", err)
		return
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error: read help.md file fail:%v", err)
		return
	}

	fmt.Printf(string(contents))
}*/
