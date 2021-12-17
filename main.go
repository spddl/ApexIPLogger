package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MakeNowJust/hotkey"
	"github.com/google/gopacket/pcap"
)

var (
	currentIP     string
	queueFirewall QueueFirewall
	maxDelay      = time.Millisecond * time.Duration(100)
)

func main() {
	cleanUpChan := make(chan os.Signal, 1)
	signal.Notify(cleanUpChan, syscall.SIGTERM, syscall.SIGINT)

	queueFirewall.data = make(map[string]struct{})
	hkey := hotkey.New()

	var err error
	currentIP, err = localIP()
	if err != nil {
		log.Fatal(err)
	}
	currentDevice := getCurrentDevice(currentIP)

	done := make(chan bool)
	go StartHotkeyListener(done)

	id_p, err := hkey.Register(hotkey.Ctrl|hotkey.Shift, 'P', func() {
		packetLoop(currentDevice)
	})
	if err != nil {
		log.Fatal(err)
	}
	id_o, err := hkey.Register(hotkey.Ctrl|hotkey.Shift, 'O', func() {
		log.Printf("\033[31mServerIP: %s will be blocked\033[0m\n", lastServerIP)
		go queueFirewall.addFirewallQueue(lastServerIP)
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Ctrl+Shift+P => ip search")
	log.Println("Ctrl+Shift+O => the last ip will be blocked")

	for {
		<-cleanUpChan
		hkey.Unregister(id_p)
		hkey.Unregister(id_o)
		os.Exit(0)
	}
}

func getCurrentDevice(currentIP string) string {
	// Find all devices
	devices, err := pcap.FindAllDevs() // https://haydz.github.io/2020/07/06/Go-Windows-NIC.html
	if err != nil {
		log.Fatal(err)
	}

	// device information
	for _, device := range devices {
		for _, v := range device.Addresses {
			if v.IP.String() == currentIP {
				// log.Printf("Found: %s => %s\n", device.Description, device.Name)
				return device.Name
			}
		}
	}
	return ""
}
