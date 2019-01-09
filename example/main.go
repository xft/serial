package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/xft/serial"
)

var (
	device   string
	baudrate int
	databits int
	stopbits int
	parity   string

	message string
)

func main() {
	flag.StringVar(&device, "D", "/dev/ttyUSB0", "device")
	flag.IntVar(&baudrate, "b", 115200, "baud rate")
	flag.IntVar(&databits, "d", 8, "data bits")
	flag.IntVar(&stopbits, "s", 1, "stop bits")
	flag.StringVar(&parity, "p", "N", "parity (N/E/O)")
	flag.StringVar(&message, "m", "serial", "message")
	flag.Parse()

	c := serial.Config{
		Device:   device,
		BaudRate: baudrate,
		DataBits: byte(databits),
		StopBits: serial.StopBits(stopbits),
		Parity:   serial.Parity(parity[0]),
	}

	p, err := serial.Open(&c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("open %v success\n", c.Device)
	defer func() {
		err := p.Close()
		if err != nil {
			log.Fatalf("close %v error: %v", c.Device, err)
		}
		log.Printf("close %v success\n", c.Device)
	}()

	//err = p.Flush()
	//if err != nil {
	//	log.Fatalf("Flush error: %v", err)
	//}

	n, err := p.Write([]byte(message))
	if err != nil {
		log.Fatalf("write %v error: %v", c.Device, err)
		return
	}
	fmt.Printf("write %v success: %v bytes\n", c.Device, n)

	if _, err = io.Copy(os.Stdout, p); err != nil {
		log.Fatalf("io.Copy error: %v", err)
		return
	}
}
