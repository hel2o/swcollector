package g

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/hel2o/management-system/models"
	"github.com/hel2o/management-system/tools"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var KPR *tools.KeyPairReload

func StartSSL() {
	var err error
	KPR, err = tools.NewKeyPairReload("network.clifford.cn")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	go tools.AutoReport(models.AppReportParams{
		Agent:     "swcollector",
		Host:      "network.clifford.cn",
		Version:   VERSION,
		ReloadSSL: true,
		Journal:   true,
		Api:       "https://network.clifford.cn:1989/reloadssl",
		StartTime: time.Now(),
	})
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

func DeleteSlice(a []string, elem string) []string {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
