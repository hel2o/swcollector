package g

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/toolkits/slice"

	"time"

	netTool "github.com/toolkits/net"
)

var Root string

func InitRootDir() {
	var err error
	Root, err = os.Getwd()
	if err != nil {
		log.Fatalln("getwd fail:", err)
	}
}

var LocalIps []string
var StartTime int64
var LocalIp string

func InitLocalIp() {
	if Config().Rpc.Enabled {
		conn, err := net.DialTimeout("tcp", "192.168.99.118:80", time.Second*10)
		if err != nil {
			log.Println("get local addr failed !", err)
		} else {
			LocalIp = strings.Split(conn.LocalAddr().String(), ":")[0]
			conn.Close()
		}
	} else {
		log.Println("rpc is not enabled, can't get localip")
	}
}
func InitLocalIps() {
	var err error
	LocalIps, err = netTool.IntranetIP()
	if err != nil {
		log.Fatalln("get intranet ip fail:", err)
	}
	StartTime = time.Now().Unix()
}

var (
	ips     []string
	ipsLock = new(sync.Mutex)
)

func TrustableIps() []string {
	ipsLock.Lock()
	defer ipsLock.Unlock()
	ips := Config().Http.TrustIps
	return ips
}

func IsTrustable(remoteAddr string) bool {
	ip := remoteAddr
	idx := strings.LastIndex(remoteAddr, ":")
	if idx > 0 {
		ip = remoteAddr[0:idx]
	}

	if ip == "127.0.0.1" {
		return true
	}

	return slice.ContainsString(TrustableIps(), ip)
}
