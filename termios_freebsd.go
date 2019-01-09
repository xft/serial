package serial

import "C"
import (
	"fmt"

	"golang.org/x/sys/unix"
)

var baudRateMap = map[int]uint32{
	50:     unix.B50,
	75:     unix.B75,
	110:    unix.B110,
	134:    unix.B134,
	150:    unix.B150,
	200:    unix.B200,
	300:    unix.B300,
	600:    unix.B600,
	1200:   unix.B1200,
	1800:   unix.B1800,
	2400:   unix.B2400,
	4800:   unix.B4800,
	9600:   unix.B9600,
	19200:  unix.B19200,
	38400:  unix.B38400,
	57600:  unix.B57600,
	115200: unix.B115200,
	230400: unix.B230400,
	460800: unix.B460800,
	921600: unix.B921600,
}

func cfsetspeed(termios *unix.Termios, baudRate int) error {
	speed, ok := baudRateMap[baudRate]
	if !ok {
		return fmt.Errorf("unsupported baud rate %v", baudRate)
	}

	termios.Ispeed = speed
	termios.Ospeed = speed

	return nil
}
