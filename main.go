package main

import (
	"net/netip"
	"syscall"
)

func printIP(ip [4]byte) string {
	return netip.AddrFrom4(ip).String()
}

func main() {

	r0, _, e1 := syscall.RawSyscall(syscall.SYS_SOCKET,
		uintptr(syscall.AF_INET),
		uintptr(syscall.SOCK_STREAM),
		uintptr(0))
	sock := int(r0)
	if e1 != 0 {
		panic(e1)
	}
	err := syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		panic(err)
	}
	s := &syscall.SockaddrInet4{
		Port: 6969,
		Addr: [4]byte{0x7f, 0x000, 0x00, 0x01}}
	if err := syscall.Bind(sock, s); err != nil {
		panic(err)
	}
	_, _, e1 = syscall.Syscall(syscall.SYS_LISTEN, uintptr(sock), uintptr(1), 0)
	if e1 != 0 {
		panic(e1)
	}
	buf := make([]byte, 1024)
	println("Listening on", printIP(s.Addr), s.Port)
	for {
		conn, a, err := syscall.Accept(sock)
		if err != nil {
			panic(err)
		}
		c := (a).(*syscall.SockaddrInet4)
		println("Client connected:", printIP(c.Addr), c.Port)
		if err = syscall.Sendto(conn, []byte("Welcome!\r\n"), 0, nil); err != nil {
			panic(err)
		}
		for {
			n, _, err := syscall.Recvfrom(conn, buf, 0)
			if err != nil {
				panic(err)
			}
			if n == 0 {
				println("Client disconnected:", printIP(c.Addr), c.Port)
				break
			}
			data := buf[:n]
			println("Received:", string(data))
			if err = syscall.Sendto(conn, data, 0, nil); err != nil {
				panic(err)
			}
			clear(buf)
		}
	}
}
