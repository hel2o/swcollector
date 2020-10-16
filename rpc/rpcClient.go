package rpc

import (
	"errors"
	"log"
	"math"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/net"
)

type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	RpcServer string
	Timeout   time.Duration
}

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}
}

func (this *SingleConnRpcClient) insureConn() {
	if this.rpcClient != nil {
		return
	}

	var err error
	var retry int = 1

	for {
		if this.rpcClient != nil {
			return
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.RpcServer, this.Timeout)
		if err == nil {
			return
		}

		log.Printf("dial %s fail: %v", this.RpcServer, err)

		if retry > 6 {
			retry = 1
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)

		retry++
	}
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	this.insureConn()

	timeout := time.Duration(50 * time.Second)
	done := make(chan error)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("[WARN] rpc call timeout %v => %v", this.rpcClient, this.RpcServer)
		this.close()
		return errors.New(this.RpcServer + " rpc call timeout")
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}

var (
	TransferClient *SingleConnRpcClient
)

func InitRpcClients() {
	if g.Config().Transfer.Enabled {
		TransferClient = &SingleConnRpcClient{
			RpcServer: g.Config().Transfer.Addr,
			Timeout:   time.Duration(g.Config().Transfer.Timeout) * time.Millisecond,
		}
	}
}

func SendToTransfer(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := g.Config().Debug
	debug_endpoints := g.Config().Debugmetric.Endpoints
	debug_metrics := g.Config().Debugmetric.Metrics
	debug_tags := g.Config().Debugmetric.Tags
	debug_Tags := strings.Split(debug_tags, ",")

	if g.Config().SwitchHosts.Enabled {
		hosts := g.HostConfig().Hosts
		for i, metric := range metrics {
			if hostname, ok := hosts[metric.Endpoint]; ok {
				metrics[i].Endpoint = hostname
			}
		}
	}

	if debug {
		for _, metric := range metrics {
			metric_tags := strings.Split(metric.Tags, ",")
			if g.In_array(metric.Endpoint, debug_endpoints) && g.In_array(metric.Metric, debug_metrics) {
				if debug_tags == "" {
					log.Printf("=> <Total=%d> %v\n", len(metrics), metric)
					continue
				}
				if g.Array_include(debug_Tags, metric_tags) {
					log.Printf("=> <Total=%d> %v\n", len(metrics), metric)
				}
			}
		}
	}
	var resp model.TransferResponse
	err := TransferClient.Call("Transfer.Update", metrics, &resp)
	if err != nil {
		log.Println("call Transfer.Update fail", err)
		if debug {
			for _, metric := range metrics {
				log.Printf("=> <Total=%d> %v\n", len(metrics), metric)
			}
		}
	}
	log.Println("<=", &resp)
}
