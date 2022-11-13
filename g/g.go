package g

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func ModifyRlimit() {
	var rLimit syscall.Rlimit
	rLimit.Max = 999999
	rLimit.Cur = 999999
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatal("Error Setting Rlimit ", err)
	}
}

func GetCommunity(ip string) (community string) {
	community = Config().Switch.Community
	for _, ipc := range Config().Switch.IpRange {
		if ip == ipc.Ip {
			community = ipc.Community
			break
		}
	}
	return
}
func GetHostname(ip string) (hostname string) {
	for _, ipc := range Config().Switch.IpRange {
		if ip == ipc.Ip {
			hostname = ipc.Hostname
			break
		}
	}
	return
}
func InArray(str string, array []string) bool {
	for _, s := range array {
		if str == s {
			return true
		}
	}
	return false
}

var (
	logFileName = flag.String("log", "var/app.log", "Log file name")
)

func MyLog() {
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "APP start Failed")
		os.Exit(1)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
