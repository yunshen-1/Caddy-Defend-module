package main

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"testing"
)

// ipToAddrFast is a function to convert net.IP to netip.Addr with potential optimizations
func ipToAddrFast(ip net.IP) (netip.Addr, error) { /*Slowest*/
	if ip == nil {
		return netip.Addr{}, fmt.Errorf("ip is nil")
	}
	if ip4 := ip.To4(); ip4 != nil {
		var addr4 [4]byte
		copy(addr4[:], ip4)
		return netip.AddrFrom4(addr4), nil
	}
	if ip16 := ip.To16(); ip16 != nil {
		var addr16 [16]byte
		copy(addr16[:], ip16)
		return netip.AddrFrom16(addr16), nil
	}
	return netip.Addr{}, fmt.Errorf("invalid IP address")
}

// ipToAddr is the function provided by the user for comparison
func netIPToNetipAddr(ip net.IP) (netip.Addr, error) {
	if ip4 := ip.To4(); ip4 != nil {
		return netip.AddrFrom4([4]byte{ip4[0], ip4[1], ip4[2], ip4[3]}), nil
	}

	if len(ip) == net.IPv6len {
		var addr [16]byte
		copy(addr[:], ip)
		return netip.AddrFrom16(addr), nil
	}

	return netip.Addr{}, errors.New("invalid IP address format")
}

func ipToAddrStd(ip net.IP) (netip.Addr, error) {
	if ip == nil {
		return netip.Addr{}, fmt.Errorf("ip is nil")
	}

	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return netip.Addr{}, fmt.Errorf("invalid IP address")
	}
	return addr, nil
}

//
//func main() {
//	ip4 := net.ParseIP("192.168.1.1")
//	ip6 := net.ParseIP("2001:db8::1")
//	invalidIP := net.IP([]byte{1, 2, 3})
//
//	addr4Fast, errFast := ipToAddrFast(ip4)
//	addr6Fast, errFast := ipToAddrFast(ip6)
//	_, errFastInvalid := ipToAddrFast(invalidIP)
//	if errFastInvalid == nil {
//		panic("expected error for invalid IP")
//	}
//
//	addr4User, errUser := netIPToNetipAddr(ip4)
//	addr6User, errUser := netIPToNetipAddr(ip6)
//	_, errUserInvalid := netIPToNetipAddr(invalidIP)
//	if errUserInvalid == nil {
//		panic("expected error for invalid IP")
//	}
//
//	addr4Std, errStd := ipToAddrStd(ip4)
//	addr6Std, errStd := ipToAddrStd(ip6)
//	_, errStdInvalid := ipToAddrStd(invalidIP)
//	if errStdInvalid == nil {
//		panic("expected error for invalid IP")
//	}
//
//	fmt.Println("Fast IPv4:", addr4Fast)
//	fmt.Println("Fast IPv6:", addr6Fast)
//	fmt.Println("User IPv4:", addr4User)
//	fmt.Println("User IPv6:", addr6User)
//	fmt.Println("Std IPv4:", addr4Std)
//	fmt.Println("Std IPv6:", addr6Std)
//
//	if errFast != nil || errUser != nil || errStd != nil {
//		panic("unexpected error during conversion")
//	}
//}

func BenchmarkIPToAddrFastIPv4(b *testing.B) {
	ip := net.ParseIP("192.168.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ipToAddrFast(ip)
	}
}

func BenchmarkNetIPToNetipAddrIPv4(b *testing.B) {
	ip := net.ParseIP("192.168.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = netIPToNetipAddr(ip)
	}
}

func BenchmarkIPToAddrStdIPv4(b *testing.B) {
	ip := net.ParseIP("192.168.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ipToAddrStd(ip)
	}
}

func BenchmarkIPToAddrFastIPv6(b *testing.B) {
	ip := net.ParseIP("2001:db8::1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ipToAddrFast(ip)
	}
}

func BenchmarkNetIPToNetipAddrIPv6(b *testing.B) {
	ip := net.ParseIP("2001:db8::1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = netIPToNetipAddr(ip)
	}
}

func BenchmarkIPToAddrStdIPv6(b *testing.B) {
	ip := net.ParseIP("2001:db8::1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ipToAddrStd(ip)
	}
}
