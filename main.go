package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Addr struct {
	network string
	address string
}

func (addr Addr) String() string {
	return addr.network + "://" + addr.address
}

type Config struct {
	listen Addr
	target Addr
}

func main() {
	log.Println("Start Local Port Forwarding")

	config, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listen: '%s', target: '%s'", config.listen.String(), config.target.String())

	listener, err := net.Listen(config.listen.network, config.listen.address)
	if err != nil {
		log.Fatal(err)
	}

	if err := config.testTarget(); err != nil {
		log.Fatal(err)
	}

	config.serve(listener)

	log.Println("End Local Port Forwarding")
}

func parseConfig() (*Config, error) {
	listen := flag.String("listen", "", "listen address")
	target := flag.String("target", "", "target address")
	flag.Parse()

	listenAddress, err := parseAddress(*listen)
	if err != nil {
		return nil, fmt.Errorf("error when parsing listen address: %w", err)
	}

	targetAddress, err := parseAddress(*target)
	if err != nil {
		return nil, fmt.Errorf("error when parsing target address: %w", err)
	}

	cfg := &Config{
		listen: listenAddress,
		target: targetAddress,
	}
	return cfg, nil
}

func parseAddress(address string) (Addr, error) {
	address = strings.Trim(address, " ")
	if address == "" {
		return Addr{"", ""}, errors.New("address cannot be empty")
	}

	if strings.HasPrefix(address, ":") && !strings.HasPrefix(address, "://") {
		return Addr{"tcp", "localhost" + address}, nil
	}

	addrParts := strings.Split(address, "://")
	if len(addrParts) > 2 {
		return Addr{"", ""}, errors.New("the address cannot contain more than one '://'")
	}

	if len(addrParts) == 1 {
		return Addr{"tcp", addrParts[0]}, nil
	}

	if addrParts[0] == "" {
		return Addr{"", ""}, errors.New("network not specified")
	}

	return Addr{addrParts[0], addrParts[1]}, nil
}

func (c *Config) testTarget() error {
	target, err := net.Dial(c.target.network, c.target.address)
	if err != nil {
		return err
	}
	_ = target.Close()

	return nil
}

func (c *Config) serve(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go c.handleConn(conn)
	}
}

func (c *Config) handleConn(serverConn net.Conn) {
	clientConn, err := net.Dial(c.target.network, c.target.address)
	if err != nil {
		return
	}

	cp := func(dst, src net.Conn) {
		defer func() {
			_ = dst.Close()
			_ = src.Close()
		}()

		_, _ = io.Copy(dst, src)
	}

	go cp(serverConn, clientConn)
	cp(clientConn, serverConn)
}
