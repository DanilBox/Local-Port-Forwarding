package main

import "testing"

func TestParseAddressPositive(t *testing.T) {
	testsPositive := []struct {
		address string
		want    string
	}{
		{
			"tcp://localhost:8080",
			"tcp://localhost:8080",
		},
		{
			"unix://var/lib/socket.sock",
			"unix://var/lib/socket.sock",
		},
		{
			"tcp4://user:pass@127.0.0.1:9000",
			"tcp4://user:pass@127.0.0.1:9000",
		},
		{
			"user:pass@127.0.0.1:8123",
			"tcp://user:pass@127.0.0.1:8123",
		},
		{
			"udp6://srv1525:3306",
			"udp6://srv1525:3306",
		},
		{
			"localhost:8080",
			"tcp://localhost:8080",
		},
		{
			":8080",
			"tcp://localhost:8080",
		},
	}

	for _, tt := range testsPositive {
		t.Run(tt.address, func(t *testing.T) {
			addr, err := parseAddress(tt.address)
			if err != nil {
				t.Errorf("want: 'no err', got err: '%v'", err)
			}

			got := addr.String()
			if tt.want != got {
				t.Errorf("want: '%s', got: '%s'", tt.want, got)
			}
		})
	}
}

func TestParseAddressNegative(t *testing.T) {
	testsNegative := []struct {
		address string
		want    string
	}{
		{
			"",
			"address cannot be empty",
		},
		{
			"tcp://localhost://localhost",
			"the address cannot contain more than one '://'",
		},
		{
			"://localhost",
			"network not specified",
		},
	}

	for _, tt := range testsNegative {
		t.Run(tt.address, func(t *testing.T) {
			addr, err := parseAddress(tt.address)

			if err == nil {
				got := addr.String()
				t.Errorf("want: 'err', got: '%s'", got)
			} else if tt.want != err.Error() {
				t.Errorf("want: '%s', got: '%s'", tt.want, err.Error())
			}
		})
	}
}
