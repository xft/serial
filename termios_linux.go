package serial

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

var baudRateMap = map[int]uint32{
	50:      unix.B50,
	75:      unix.B75,
	110:     unix.B110,
	134:     unix.B134,
	150:     unix.B150,
	200:     unix.B200,
	300:     unix.B300,
	600:     unix.B600,
	1200:    unix.B1200,
	1800:    unix.B1800,
	2400:    unix.B2400,
	4800:    unix.B4800,
	9600:    unix.B9600,
	19200:   unix.B19200,
	38400:   unix.B38400,
	57600:   unix.B57600,
	115200:  unix.B115200,
	230400:  unix.B230400,
	460800:  unix.B460800,
	500000:  unix.B500000,
	576000:  unix.B576000,
	921600:  unix.B921600,
	1000000: unix.B1000000,
	1152000: unix.B1152000,
	1500000: unix.B1500000,
	2000000: unix.B2000000,
	2500000: unix.B2500000,
	3000000: unix.B3000000,
	3500000: unix.B3500000,
	4000000: unix.B4000000,
}

func cfsetspeed(termios *unix.Termios, baudRate int) error {
	speed, ok := baudRateMap[baudRate]
	if !ok {
		return fmt.Errorf("unsupported baud rate %v", baudRate)
	}

	termios.Cflag &^= unix.CBAUD | unix.CBAUDEX
	termios.Cflag |= speed
	termios.Ispeed = speed
	termios.Ospeed = speed

	return nil
}

func tcgetattr(file *os.File) (termios *unix.Termios, err error) {
	termios = &unix.Termios{}
	if _, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		file.Fd(),
		uintptr(unix.TCGETS),
		uintptr(unsafe.Pointer(termios)),
	); errno != 0 {
		return nil, fmt.Errorf("tcgetattr: %v", errno)
	}
	return
}

func tcsetattr(file *os.File, optionalActions int, termios *unix.Termios) (err error) {
	cmd := unix.TCSETS
	switch optionalActions {
	case tioTCSANOW:
		cmd = unix.TCSETS
	case tioTCSADRAIN:
		cmd = unix.TCSETSW
	case tioTCSAFLUSH:
		cmd = unix.TCSETSF
	}
	if _, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		file.Fd(),
		uintptr(cmd),
		uintptr(unsafe.Pointer(termios)),
	); errno != 0 {
		return fmt.Errorf("tcsetattr: %v", errno)
	}
	return
}

func tcflush(file *os.File) error {
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(unix.TCFLSH),
		uintptr(unix.TCIOFLUSH),
	)
	if errno == 0 {
		return nil
	}
	return errno
}
