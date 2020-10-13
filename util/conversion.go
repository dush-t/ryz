package util

import (
	"encoding/binary"
	"fmt"
	"net"
)

func IpToBinary(ipStr string) ([]byte, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("Not a valid IP: %s", ipStr)
	}
	return []byte(ip.To4()), nil
}

func MacToBinary(macStr string) ([]byte, error) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, err
	}
	return []byte(mac), nil
}

func UInt32ToBinary(i uint32, numBytes int) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b[numBytes:], nil
}

// I know, alright?
func BinaryToUint32(data []byte) uint32 {
	return uint32(uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16 + uint32(data[3])<<24)
}

// I couldn't find another way.
func Binary48ToInt64(data []byte) uint64 {
	return uint64(uint64(data[0]) + uint64(data[1])<<8 + uint64(data[2])<<16 + uint64(data[3])<<24 + uint64(data[4])<<32 + uint64(data[5])<<40)
}
