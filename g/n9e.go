package g

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"strings"
	"time"

	"github.com/didi/nightingale/src/common/dataobj"

	"github.com/open-falcon/common/model"
	"github.com/ugorji/go/codec"
)

func N9ePush(items []*model.MetricValue) {
	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	addr := config.Transfer.N9e
	retry := 0
	for {
		conn, err := net.DialTimeout("tcp", addr, time.Millisecond*3000)
		if err != nil {
			log.Println("dial transfer err:", err)
			continue
		}

		var bufconn = struct { // bufconn here is a buffered io.ReadWriteCloser
			io.Closer
			*bufio.Reader
			*bufio.Writer
		}{conn, bufio.NewReader(conn), bufio.NewWriter(conn)}

		rpcCodec := codec.MsgpackSpecRpc.ClientCodec(bufconn, &mh)
		client := rpc.NewClientWithCodec(rpcCodec)

		debug := Config().Debug
		debug_endpoints := Config().Debugmetric.Endpoints
		debug_items := Config().Debugmetric.Metrics
		debug_tags := Config().Debugmetric.Tags
		debug_Tags := strings.Split(debug_tags, ",")

		if Config().SwitchHosts.Enabled {
			hosts := HostConfig().Hosts
			for i, metric := range items {
				if hostname, ok := hosts[metric.Endpoint]; ok {
					items[i].Endpoint = hostname
				}
			}
		}

		if debug {
			for _, metric := range items {
				metric_tags := strings.Split(metric.Tags, ",")
				if In_array(metric.Endpoint, debug_endpoints) && In_array(metric.Metric, debug_items) {
					if debug_tags == "" {
						log.Printf("=> <Total=%d> %v\n", len(items), metric)
						continue
					}
					if Array_include(debug_Tags, metric_tags) {
						log.Printf("=> <Total=%d> %v\n", len(items), metric)
					}
				}
			}
		}

		var reply dataobj.TransferResp
		err = client.Call("Transfer.Push", items, &reply)
		client.Close()
		if err != nil {
			log.Println("发送数据到N9E出错", err)
			continue
		} else {
			log.Println("N9E回复", reply.String())
			return
		}
		time.Sleep(time.Millisecond * 500)

		retry += 1
		if retry == 3 {
			retry = 0
			break
		}
	}
}

func Array_include(array_a []string, array_b []string) bool { //b include a
	for _, v := range array_a {
		if In_array(v, array_b) {
			continue
		} else {
			return false
		}
	}
	return true
}

func In_array(a string, array []string) bool {
	for _, v := range array {
		if a == v {
			return true
		}
	}
	return false
}
