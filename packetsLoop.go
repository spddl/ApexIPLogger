package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var runningPacketLoop = false
var lastServerIP string

func packetLoop(currentDevice string) {
	if runningPacketLoop {
		return
	}
	runningPacketLoop = true
	if handle, err := pcap.OpenLive(currentDevice, 1600, true, pcap.BlockForever); err != nil {
		log.Fatal(err)
	} else if err := handle.SetBPFFilter("udp and portrange 37005-38515"); err != nil {
		log.Fatal(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for {
			packet, err := packetSource.NextPacket()
			if err != nil {
				continue
			}

			serverIP := printPacketInfo(packet)
			if serverIP == "" {
				continue
			}

			_, t, err := Ping(serverIP)
			if err != nil {
				log.Println(err) // read ip4 0.0.0.0: i/o timeout
				continue
			}

			lastServerIP = serverIP
			if t <= maxDelay {
				// under 100ms
				log.Printf("\033[32mServerIP: %s, %v\033[0m\n", serverIP, t)
			} else {
				// over 100ms
				log.Printf("\033[31mServerIP: %s will be blocked, %v\033[0m\n", serverIP, t)
				go queueFirewall.addFirewallQueue(serverIP)
			}
			runningPacketLoop = false
			handle.Close()
			break
		}

	}
}

func printPacketInfo(packet gopacket.Packet) string {
	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)

		srcIP := ip.SrcIP.String()

		var dest_ip string
		if currentIP == srcIP {
			dest_ip = ip.DstIP.String()
		} else {
			dest_ip = srcIP
		}

		return dest_ip
	}
	return ""
}
