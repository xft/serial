package serial

import "C"
import (
	"fmt"

	"golang.org/x/sys/unix"
)

func cfsetspeed(termios *unix.Termios, baudRate int) error {
	speed := termios.Ispeed

	switch baudRate {
	case 0:
		speed = unix.B115200
	case 50:
		speed = unix.B50
	case 75:
		speed = unix.B75
	case 110:
		speed = unix.B110
	case 134:
		speed = unix.B134
	case 150:
		speed = unix.B150
	case 200:
		speed = unix.B200
	case 300:
		speed = unix.B300
	case 600:
		speed = unix.B600
	case 1200:
		speed = unix.B1200
	case 1800:
		speed = unix.B1800
	case 2400:
		speed = unix.B2400
	case 4800:
		speed = unix.B4800
	case 7200:
		speed = unix.B7200
	case 9600:
		speed = unix.B9600
	case 14400:
		speed = unix.B14400
	case 19200:
		speed = unix.B19200
	case 28800:
		speed = unix.B28800
	case 38400:
		speed = unix.B38400
	case 57600:
		speed = unix.B57600
	case 76800:
		speed = unix.B76800
	case 115200:
		speed = unix.B115200
	case 230400:
		speed = unix.B230400
	default:
		return fmt.Errorf("unsupported baud rate %v", baudRate)
	}

	termios.Ispeed = speed
	termios.Ospeed = speed

	return nil
}
