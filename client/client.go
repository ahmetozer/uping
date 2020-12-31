package client

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	share "github.com/ahmetozer/uping/share"
	"github.com/beevik/ntp"
)

// Main Client main function
func Main(args []string) {
	fmt.Println("Client Function Executed")
	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
	timeServer := clientCmd.String("timeserver", "time.cloudflare.com", "Time server which is used for updating time")
	timeSyncInterval := clientCmd.Uint("tsi", 10, "Time sync interval for time client")
	pingInterval := clientCmd.Float64("i", 1, "Ping interval")
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

	if err := share.NtpUpdate(&timeOffset, *timeServer); err != nil {
		log.Fatalf("First time sync is not run %v", err)
	}

	fmt.Printf("Time synced, the offset is \"%v\"\n", timeOffset.ClockOffset)

	go share.NtpUpdateLoop(&timeOffset, *timeServer, *timeSyncInterval)

	if share.IsMayOnlyIPv6(remoteAddr) {
		remoteAddr = "[" + remoteAddr + "]"
	}

	p := make([]byte, 6)
	conn, err := net.Dial("udp", remoteAddr+":50123")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	defer conn.Close()
	// Sender
	for {
		timeByteArr := share.TimeIntToByteArray(int32((time.Now().Add((*timeOffset).ClockOffset).UnixNano() / 1000000) % 10000000))
		// Append same data to twice and equalize inbound and outbound byte
		conn.Write(append(timeByteArr, timeByteArr...))
		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			time3 := int32((time.Now().Add((*timeOffset).ClockOffset).UnixNano() / 1000000) % 10000000)
			time1 := share.TimeByteArrayToInt(p[0:3])
			time2 := share.TimeByteArrayToInt(p[3:6])
			fmt.Printf("6 bytes client > %vms > server(%v) > %vms > client total %v\n", time2-time1, remoteAddr, time3-time2, time3-time1)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
		time.Sleep(time.Duration(*pingInterval*1000) * time.Millisecond)
	}
}
