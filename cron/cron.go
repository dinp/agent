package cron

import (
	"github.com/dinp/agent/g"
	"github.com/dinp/common/dock"
	"github.com/dinp/common/model"
	"github.com/toolkits/nux"
	"log"
	"time"
)

func Heartbeat() {
	duration := time.Duration(g.Config().Interval) * time.Second
	for {
		heartbeat(duration)
		time.Sleep(duration)
	}
}

func heartbeat(duration time.Duration) {
	mem, err := nux.MemInfo()
	if err != nil {
		log.Println("[ERROR] get meminfo:", err)
		return
	}

	// use MB
	memFree := (mem.MemFree + mem.Buffers + mem.Cached) / 1024 / 1024

	localIp := g.Config().LocalIp

	containers, err := dock.Containers(g.Config().Docker)
	if err != nil {
		log.Println("[ERROR] list containers fail:", err)

		// this node dead
		var resp model.NodeResponse
		err = g.RpcClient.Call("NodeState.NodeDown", localIp, &resp)
		if err != nil || resp.Code != 0 {
			log.Println("[ERROR] call rpc: NodeState.NodeDown fail:", err, "resp:", resp)
		} else if g.Config().Debug {
			log.Println("[INFO] call rpc: NodeState.NodeDown successfully. I am dead...")
		}

		for {
			time.Sleep(duration)
			containers, err = dock.Containers(g.Config().Docker)
			if err == nil {
				break
			} else {
				log.Println("[ERROR] list containers fail:", err)
			}
		}
	}

	req := model.NodeRequest{
		Node: model.Node{
			Ip:      localIp,
			MemFree: memFree,
		},
		Containers: containers,
	}
	var resp model.NodeResponse
	err = g.RpcClient.Call("NodeState.Push", req, &resp)
	if err != nil || resp.Code != 0 {
		log.Println("[ERROR] call rpc: NodeState.Push fail:", err, "resp:", resp)
	} else if g.Config().Debug {
		log.Println("[DEBUG] =>", req)
	}
}
