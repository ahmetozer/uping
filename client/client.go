package client

import (
	"bufio"
	"flag"
	"fmt"
	share "github.com/ahmetozer/uping/share"
	"github.com/beevik/ntp"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Main Client main function
func Main(args []string) {
	fmt.Println("Client Mode")
	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
	timeServer := clientCmd.String("timeserver", "time.cloudflare.com", "Time server which is used for updating time")
	timeSyncInterval := clientCmd.Uint("tsi", 10, "Time sync interval for time client")
	pingInterval := clientCmd.Float64("i", 1, "Ping interval")
	remotePort := clientCmd.Uint("p", 50123, "Remote port")
	sourcePort := clientCmd.Uint("sp", 0, "Source port")
	pingCount := clientCmd.Uint("c", 0, "Ping count")
	clientCmd.Parse(args)

	if *timeSyncInterval < 10 {
		log.Fatalf("time sync update interval \"%v\" is to low\n", *timeSyncInterval)
	}

	if !share.IsMayHost(*timeServer) {
		log.Fatalf("TimeServer \"%v\" is not a host!", timeServer)
	}

	if *pingInterval < 0.1 {
		log.Printf("Ping interval \"%v\" is too low, setted to lowest value \"0.1\"", *pingInterval)
		*pingInterval = 0.1
	}

	clientCmdTail := clientCmd.Args()
	if len(clientCmdTail) == 0 {
		log.Fatalf("Remote address expected")
	} else if len(clientCmdTail) != 1 {
		log.Fatalf("Only one remote address required")
	}
	remoteAddr := clientCmdTail[0]
	if !share.IsMayHost(remoteAddr) {
		log.Fatalf("Remote address \"%v\" is not a host !", remoteAddr)
	}

	var timeOffset *ntp.Response
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	if err := share.NtpUpdate(&timeOffset, *timeServer); err != nil {
		log.Fatalf("First time sync is not run %v", err)
	}

	var sendedPacket uint32
	var receivedPacket uint32
	var avgOutgoingLatency int64
	var avgIncomingLatency int64
	var killListen uint

	go func() {
		sig := <-signals
		fmt.Println(sig)
		printStatistic(sendedPacket, receivedPacket, avgOutgoingLatency, avgIncomingLatency)
		os.Exit(0)
	}()

	fmt.Printf("Time synced, the offset is \"%v\"\n", timeOffset.ClockOffset)

	go share.NtpUpdateLoop(&timeOffset, *timeServer, *timeSyncInterval)

	if share.IsMayOnlyIPv6(remoteAddr) {
		remoteAddr = "[" + remoteAddr + "]"
	}

	p := make([]byte, 6)
	sp, err := strconv.Atoi(fmt.Sprintf("%v",*sourcePort));
	if err != nil {
		// handle error
		log.Fatalf("Source port err: %v",err)
		os.Exit(2)
	}

	var dialer = &net.Dialer{
		LocalAddr: &net.UDPAddr{
			Port: sp,
		},
	}
	conn, err := dialer.Dial("udp", fmt.Sprintf("%v:%v", remoteAddr, *remotePort))
	if err != nil {
		log.Fatalf("error %v\n", err)
		return
	}
	fmt.Printf("\nPinging from %v to %v:%v\n", conn.LocalAddr(), remoteAddr, *remotePort)

	defer conn.Close()

	// Sender
	go func() {
		for {
			timeByteArr := share.TimeIntToByteArray(int32((time.Now().Add((*timeOffset).ClockOffset).UnixNano() / 1000000) % 10000000))
			conn.Write(append(timeByteArr, timeByteArr...))
			sendedPacket = sendedPacket + 1
			if *pingCount > 0 {
				*pingCount = *pingCount - 1
				if *pingCount == 0 {
					conn.SetReadDeadline(time.Now().Add(1 * time.Second))
					killListen = 1
					break
				}
			}
			time.Sleep(time.Duration(*pingInterval*1000) * time.Millisecond)
		}
	}()

	//avgIncomingLatency = -1
	avgOutgoingLatency = -1
	avgIncomingLatency = -1
	for {
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			time3 := int32((time.Now().Add((*timeOffset).ClockOffset).UnixNano() / 1000000) % 10000000)
			time1 := share.TimeByteArrayToInt(p[0:3])
			time2 := share.TimeByteArrayToInt(p[3:6])
			receivedPacket = receivedPacket + 1
			if avgOutgoingLatency == -1 {
				avgOutgoingLatency = int64(time2 - time1)
				avgIncomingLatency = int64(time3 - time2)
			} else {
				tempLat := (avgOutgoingLatency * int64(receivedPacket-1)) + int64((time2 - time1))
				avgOutgoingLatency = tempLat / int64(receivedPacket)
				tempLat = (avgIncomingLatency * int64(receivedPacket-1)) + int64((time3 - time2))
				avgIncomingLatency = tempLat / int64(receivedPacket)
			}
			fmt.Printf("6 bytes client > %v ms > server(%v) > %v ms > client total %v ms seq %v\n", time2-time1, remoteAddr, time3-time2, time3-time1, sendedPacket)
		} else {
			log.Printf("error %v\n", err)
		}
		time.Sleep(time.Duration(*pingInterval*1000) * time.Millisecond)
		if killListen == 2 {
			break
		}
		if killListen == 1 {
			killListen = killListen + 1
		}

	}
	printStatistic(sendedPacket, receivedPacket, avgOutgoingLatency, avgIncomingLatency)
}

func printStatistic(sendedPacket uint32, receivedPacket uint32, avgOutgoingLatency int64, avgIncomingLatency int64) {
	packetLoss := float64(receivedPacket) / float64(sendedPacket)
	packetLoss = 1 - packetLoss
	packetLoss = packetLoss * 100
	fmt.Printf("\n%v packets transmitted, %v received, %v%% packet Loss\nAverage outbound packet delay %v\nAverage inbound packet delay %v\n", sendedPacket, receivedPacket, packetLoss, avgOutgoingLatency, avgIncomingLatency)
}
