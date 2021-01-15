package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/matti/betterio"
)

func handle(ctx context.Context, conn net.Conn, remoteAddress string) {
	dialer := net.Dialer{
		Timeout: 1 * time.Second,
	}

	if upstream, err := dialer.Dial("tcp", remoteAddress); err != nil {
		log.Println(conn.RemoteAddr().String(), "dial err to", remoteAddress)
	} else {
		defer upstream.Close()

		betterio.CopyBidirUntilCloseAndReturnBytesWritten(conn, upstream)
	}
}

func proxy(ctx context.Context, localAddress string, remoteAddress string) {
	var ln net.Listener
	var err error

	if ln, err = net.Listen("tcp", localAddress); err != nil {
		log.Fatalln(err)
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		var conn net.Conn
		if conn, err = ln.Accept(); err != nil {
			return
		}

		go func() {
			defer conn.Close()
			handle(ctx, conn, remoteAddress)
		}()
	}
}
func usage() {
	fmt.Println("usage")
	os.Exit(1)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	go func(cancel context.CancelFunc) {
		s := <-sigs
		log.Println("got signal", s.String())
		cancel()
	}(cancel)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	for _, address := range os.Args[1:] {
		iface := "127.0.0.1"
		var localAddress string
		var remoteAddress string

		parts := strings.Split(address, ":")

		switch len(parts) {
		case 2:
			localAddress = net.JoinHostPort(iface, parts[1])
			remoteAddress = net.JoinHostPort(parts[0], parts[1])
		case 3:
			localAddress = net.JoinHostPort(iface, parts[0])
			remoteAddress = net.JoinHostPort(parts[1], parts[2])
		case 4:
			localAddress = net.JoinHostPort(parts[0], parts[1])
			remoteAddress = net.JoinHostPort(parts[2], parts[3])
		default:
			usage()
		}

		wg.Add(1)
		go func(localAddress string, remoteAddress string) {
			proxy(ctx, localAddress, remoteAddress)
			log.Println("proxy returned")
			wg.Done()
		}(localAddress, remoteAddress)

		log.Println(localAddress, "->", remoteAddress)
	}

	wg.Wait()
	log.Println("bye")
}
