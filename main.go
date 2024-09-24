package main

import (
	"net/netip"
	"syscall"
)

var (
	host = [4]byte{0x7f, 0x000, 0x00, 0x01}
	port = 6969
)

func ipstr(ip [4]byte) string {
	return netip.AddrFrom4(ip).String()
}

func startServer(sa *syscall.SockaddrInet4) error {
	r0, _, e1 := syscall.RawSyscall(syscall.SYS_SOCKET,
		uintptr(syscall.AF_INET),
		uintptr(syscall.SOCK_STREAM),
		uintptr(0))
	sock := int(r0)
	if e1 != 0 {
		return e1
	}
	err := syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return err
	}
	if err := syscall.Bind(sock, sa); err != nil {
		return err
	}
	_, _, e1 = syscall.Syscall(syscall.SYS_LISTEN, uintptr(sock), uintptr(5), 0)
	if e1 != 0 {
		return e1
	}
	buf := make([]byte, 1024)
	println("Listening on", ipstr(sa.Addr), sa.Port)
	for {
		conn, a, err := syscall.Accept(sock)
		if err != nil {
			return err
		}
		r0, _, e1 := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
		if e1 != 0 {
			return e1
		}
		id := int(r0)
		if id == 0 {
			c := (a).(*syscall.SockaddrInet4)
			println(ipstr(c.Addr), c.Port, "Client connected")
			if err = syscall.Sendto(conn, []byte("Welcome!\r\n"), 0, nil); err != nil {
				return err
			}
			for {
				n, _, err := syscall.Recvfrom(conn, buf, 0)
				if err != nil {
					return err
				}
				if n == 0 {
					println(ipstr(c.Addr), c.Port, "Client disconnected")
					if err := syscall.Close(conn); err != nil {
						return err
					}
					break
				}
				data := buf[:n]
				println(ipstr(c.Addr), c.Port, "Received:", string(data))
				if err = syscall.Sendto(conn, data, 0, nil); err != nil {
					return err
				}
			}
		} else {
			if err := syscall.Close(conn); err != nil {
				return err
			}
		}
	}
}

func main() {
	if err := startServer(&syscall.SockaddrInet4{Port: port, Addr: host}); err != nil {
		println(err.Error())
		syscall.Exit(1)
	}
}
