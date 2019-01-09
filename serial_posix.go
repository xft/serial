// +build linux darwin freebsd

package serial

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const (
	tioTCSANOW   = 0 // make change immediate
	tioTCSADRAIN = 1 // drain output, then change
	tioTCSAFLUSH = 2 // drain output, flush input
)

func open(device string, baudRate int, dataBits byte, parity Parity, stopBits StopBits, readTimeout time.Duration) (p *Port, err error) {
	tiosToUse, err := newTermios(baudRate, dataBits, stopBits, parity, readTimeout)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(device, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	tiosSaved, err := tcgetattr(f)
	if err != nil {
		return nil, err
	}

	err = tcsetattr(f, tioTCSANOW, tiosToUse)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tcsetattr(f, tioTCSANOW, tiosSaved)
		}
	}()

	if err = unix.SetNonblock(int(f.Fd()), false); err != nil {
		return
	}

	return &Port{f: f, savedTermios: tiosSaved}, nil
}

type Port struct {
	f            *os.File
	savedTermios *unix.Termios
}

func (p *Port) Read(b []byte) (n int, err error) {
	return p.f.Read(b)
}

func (p *Port) Write(b []byte) (n int, err error) {
	return p.f.Write(b)
}

func (p *Port) Flush() error {
	return tcflush(p.f)
}

func (p *Port) Close() (err error) {
	if p.savedTermios != nil {
		tcsetattr(p.f, tioTCSANOW, p.savedTermios)
		p.savedTermios = nil
	}
	return p.f.Close()
}

func newTermios(baudRate int, dataBits byte, stopBits StopBits, parity Parity, readTimeout time.Duration) (*unix.Termios, error) {
	t := &unix.Termios{}
	if err := cfsetspeed(t, baudRate); err != nil {
		return nil, err
	}
	t.Cflag |= unix.CREAD | unix.CLOCAL
	t.Cflag &^= unix.CRTSCTS | unix.CSIZE | unix.CSTOPB | unix.HUPCL | unix.PARENB | unix.PARODD

	switch dataBits {
	case 5:
		t.Cflag |= unix.CS5
	case 6:
		t.Cflag |= unix.CS6
	case 7:
		t.Cflag |= unix.CS7
	case 8:
		t.Cflag |= unix.CS8
	default:
		return nil, fmt.Errorf("unsupported serial data bits: %v", dataBits)
	}

	switch stopBits {
	case Stop1:
		// default is 1 stop bit
	case Stop2:
		t.Cflag |= unix.CSTOPB
	default:
		// Don't know how to set 1.5
		return nil, fmt.Errorf("unsupported stop bits %v", stopBits)
	}

	switch parity {
	case ParityNone:
	case ParityOdd:
		t.Cflag |= unix.PARENB
		t.Cflag |= unix.PARODD
	case ParityEven:
		t.Cflag |= unix.PARENB
	default:
		return nil, fmt.Errorf("unsupported parity %c", parity)
	}

	t.Iflag = unix.IGNPAR

	vmin, vtime := posixTimeoutValues(readTimeout)
	t.Cc[unix.VMIN] = vmin
	t.Cc[unix.VTIME] = vtime

	return t, nil
}

// Converts the timeout values for Linux / POSIX systems
func posixTimeoutValues(readTimeout time.Duration) (vmin uint8, vtime uint8) {
	const MAXUINT8 = 1<<8 - 1 // 255
	// set blocking / non-blocking read
	var minBytesToRead uint8 = 1
	var readTimeoutInDeci int64
	if readTimeout > 0 {
		// EOF on zero read
		minBytesToRead = 0
		// convert timeout to deciseconds as expected by VTIME
		readTimeoutInDeci = readTimeout.Nanoseconds() / 1e6 / 100
		// capping the timeout
		if readTimeoutInDeci < 1 {
			// min possible timeout 1 Deciseconds (0.1s)
			readTimeoutInDeci = 1
		} else if readTimeoutInDeci > MAXUINT8 {
			// max possible timeout is 255 deciseconds (25.5s)
			readTimeoutInDeci = MAXUINT8
		}
	}
	return minBytesToRead, uint8(readTimeoutInDeci)
}
