package driver

import (
	"github.com/tarm/serial"
	"log"
	"testing"
)

func TestRun(t *testing.T) {
	c := &serial.Config{Name: "COM7", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])
}
