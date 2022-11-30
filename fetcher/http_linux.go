package fetcher

import (
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

func NewHttpClient() http.Client {
	dialer := &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, 255)
				if err != nil {
					log.Printf("control: %s\n", err)
					return
				}
			})
		},
	}

	if os.Geteuid() != 0 {
		dialer.Control = nil
	}

	return http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
