package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	share "github.com/ahmetozer/uping/share"
	"github.com/beevik/ntp"
)

// Main Server mode main program
func Main(args []string) {
	fmt.Println("Server Mode")
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	timeServer := serverCmd.String("timeserver", "time.cloudflare.com", "Time server which is used for updating time")
	timeSyncInterval := serverCmd.Uint("tsi", 10, "Time sync interval for time client")
	listenAddr := serverCmd.String("listen", ":50123", "Listen addr for this service")
	serverCmd.Parse(args)

	if *timeSyncInterval < 10 {
		log.Fatalf("time sync update interval \"%v\" is to low\n", *timeSyncInterval)
	}

	if !share.IsMayHost(*timeServer) {
		log.Fatalf("TimeServer \"%v\" is not a host!", timeServer)
	}

	serverCmd.Args()
	if len(serverCmd.Args()) != 0 {
		log.Fatalf("Argumant is not expected %v", serverCmd.Args())
	}

	var timeOffset *ntp.Response

	if err := share.NtpUpdate(&timeOffset, *timeServer); err != nil {
		log.Fatalf("First time sync is not run %v", err)
	}

	fmt.Printf("Time synced, the offset is \"%v\"\n", timeOffset.ClockOffset)

	// Sync time offest in background
	go share.NtpUpdateLoop(&timeOffset, *timeServer, *timeSyncInterval)

	conn, err := net.ListenPacket("udp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		buf := make([]byte, 3)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}
		// create thread for sending packet
		go reply(conn, addr, buf, &timeOffset)
	}
}

func reply(conn net.PacketConn, addr net.Addr, buf []byte, timeOffset **ntp.Response) {
	l := int32((time.Now().Add((*timeOffset).ClockOffset).UnixNano() / 1000000) % 10000000)
	timeByteArr := share.TimeIntToByteArray(l)
	conn.WriteTo(append(buf, timeByteArr...), addr) // Outgoing packet data size will be 6 byte
}
