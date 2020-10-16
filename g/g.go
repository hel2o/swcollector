package g

import (
	"log"
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
	if InArray(ip, Config().Switch.SpecialSw.IpRange) {
		community = Config().Switch.SpecialSw.Community
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
