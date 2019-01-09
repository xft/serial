/*
Goserial is a simple go package to allow you to read and write from
the serial port as a stream of bytes.

It aims to have the same API on all platforms, including windows.  As
an added bonus, the windows package does not use cgo, so you can cross
compile for windows from another platform.  Unfortunately goinstall
does not currently let you cross compile so you will have to do it
manually:

 GOOS=windows make clean install

Currently there is very little in the way of configurability.  You can
set the baud rate.  Then you can Read(), Write(), or Close() the
connection.  Read() will block until at least one byte is returned.
Write is the same.  There is currently no exposed way to set the
timeouts, though patches are welcome.

Currently all ports are opened with 8 data bits, 1 stop bit, no
parity, no hardware flow control, and no software flow control.  This
works fine for many real devices and many faux serial devices
including usb-to-serial converters and bluetooth serial ports.

You may Read() and Write() simulantiously on the same connection (from
different goroutines).

Example usage:

  package main

  import (
      "github.com/xft/serial"
      "log"
  )

  func main() {
        c := &serial.Config{Device: "COM5", BaudRate: 115200}
        s, err := serial.Open(c)
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
        log.Print("%q", buf[:n])
  }
*/
package serial

import (
	"time"
)

type StopBits byte
type Parity byte

const (
	Stop1     StopBits = 1
	Stop1Half StopBits = 15
	Stop2     StopBits = 2
)

const (
	ParityNone  Parity = 'N'
	ParityOdd   Parity = 'O'
	ParityEven  Parity = 'E'
	ParityMark  Parity = 'M' // parity bit is always 1
	ParitySpace Parity = 'S' // parity bit is always 0
)

// Config contains the information needed to open a serial port.
//
// Currently few options are implemented, but more may be added in the
// future (patches welcome), so it is recommended that you create a
// new config addressing the fields by name rather than by order.
//
// For example:
//
//    c0 := &serial.Config{Device: "COM45", BaudRate: 115200, ReadTimeout: time.Millisecond * 500}
// or
//    c1 := new(serial.Config)
//    c1.Device = "/dev/tty.usbserial"
//    c1.BaudRate = 115200
//    c1.ReadTimeout = time.Millisecond * 500
//
type Config struct {
	// Device path (default /dev/ttyUSB0 or COM0)
	Device string
	// Baud rate (default 115200)
	BaudRate    int
	ReadTimeout time.Duration // Total timeout

	// Number of data bits. (default 8)
	DataBits byte

	// Parity is the bit to use and defaults to ParityNone (no parity bit).
	Parity Parity

	// Number of stop bits to use. Default is 1 (1 stop bit).
	StopBits StopBits

	// RTSFlowControl bool
	// DTRFlowControl bool
	// XONFlowControl bool

	// CRLFTranslate bool
}

// Open opens a serial port with the specified configuration
func Open(c *Config) (*Port, error) {
	size, par, stop := c.DataBits, c.Parity, c.StopBits
	if size == 0 {
		size = 8
	}
	if par == 0 {
		par = ParityNone
	}
	if stop == 0 {
		stop = Stop1
	}
	return open(c.Device, c.BaudRate, size, par, stop, c.ReadTimeout)
}

// func SendBreak()

// func RegisterBreakHandler(func())
