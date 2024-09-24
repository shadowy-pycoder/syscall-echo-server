package main

import (
	"syscall"
	"testing"
)

func BenchmarkClient(b *testing.B) {
	addr := &syscall.SockaddrInet4{Port: port, Addr: host}
	buf := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		r0, _, e1 := syscall.RawSyscall(syscall.SYS_SOCKET,
			uintptr(syscall.AF_INET),
			uintptr(syscall.SOCK_STREAM),
			uintptr(0))
		sock := int(r0)
		if e1 != 0 {
			b.Fatal(e1)
		}
		if err := syscall.Connect(sock, addr); err != nil {
			b.Fatal(err)
		}
		if err := syscall.Sendto(sock, []byte("Hello World!\r\n"), 0, nil); err != nil {
			println(err.Error())
			syscall.Exit(1)
		}
		_, _, err := syscall.Recvfrom(sock, buf, 0)
		if err != nil {
			println(err.Error())
			syscall.Exit(1)
		}
		syscall.Close(sock)
	}
}
