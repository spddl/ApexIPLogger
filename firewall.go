package main

import (
	"log"
	"os/exec"
	"sync"
	"time"
)

// Open Firewall => wf.msc

type QueueFirewall struct {
	sync.RWMutex
	data map[string]struct{}
}

func setFirewallRule(ip string) error {
	_, err := exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=APEX BLOCK IP ADDRESS - "+ip, "dir=in", "action=block", "protocol=udp", "remoteip="+ip).CombinedOutput()
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(string(o))
	return err
}

func (qf *QueueFirewall) addFirewallQueue(ip string) {
	qf.Lock()
	defer qf.Unlock()

	if _, ok := qf.data[ip]; ok {
		// IP exist in Queue
		return
	}

	qf.data[ip] = struct{}{}

	for {
		err := setFirewallRule(ip)
		if err != nil {
			// retry in 5 sec
			time.Sleep(time.Second * 5)
		} else {
			log.Printf("\033[32mServerIP: %s added to the Firewall\033[0m\n", ip)
			break
		}
	}
}
