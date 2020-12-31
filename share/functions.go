package share

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unsafe"

	"github.com/beevik/ntp"
)

var (
	domainRegex = `^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z]$`
	ipv6Regex   = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	ipv4Regex   = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
)

// IsMayDomain Check is it a Domain
func IsMayDomain(host string) bool {
	match, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

// IsMayOnlyIPv4 Check is it a IPv4
func IsMayOnlyIPv4(host string) bool {
	match, err := regexp.MatchString(ipv4Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

// IsMayOnlyIPv6 Check is it a IPv6
func IsMayOnlyIPv6(host string) bool {
	match, err := regexp.MatchString(ipv6Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

// IsMayHost Check is it a Host
func IsMayHost(host string) bool {
	return IsMayOnlyIPv4(host) || IsMayOnlyIPv6(host) || IsMayDomain(host)
}

// NtpUpdate Update Computer Time
func NtpUpdate(response **ntp.Response, timeServer string) error {
	r, err := ntp.Query(timeServer)
	if err != nil {
		return fmt.Errorf("time: not synced %v", err)
	}
	*response = r
	return nil
}

// NtpUpdateLoop Ntp time loop
func NtpUpdateLoop(response **ntp.Response, timeServer string, waitTime uint) error {
	sleepDuration := time.Duration(waitTime)
	if !IsMayHost(timeServer) {
		return errors.New("TimeServer is not a host")
	}
	for {
		r, err := ntp.Query("time.cloudflare.com")
		if err != nil {
			fmt.Printf("time: not synced %v \n", err)
		} else {
			*response = r
		}
		//fmt.Printf("time: synced %v \n", r)
		time.Sleep(sleepDuration * time.Second)
	}

}

//TimeIntToByteArray unix time last ms to byte
func TimeIntToByteArray(num int32) []byte {
	//24 bit , 3 byte data
	size := int32(3)
	arr := make([]byte, size)
	for i := int32(0); i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

//TimeByteArrayToInt converting byte to unix time
func TimeByteArrayToInt(arr []byte) int32 {
	val := int32(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}
