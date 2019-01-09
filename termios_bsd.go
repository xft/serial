// +build darwin freebsd

package serial

import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func tcgetattr(file *os.File) (termios *unix.Termios, err error) {
	termios = &unix.Termios{}
	if _, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		file.Fd(),
		uintptr(unix.TIOCGETA),
		uintptr(unsafe.Pointer(termios)),
	); errno != 0 {
		return nil, fmt.Errorf("tcgetattr: %v", errno)
	}
	return
}

func tcsetattr(file *os.File, optionalActions int, termios *unix.Termios) (err error) {
	cmd := unix.TIOCSETA
	switch optionalActions {
	case tioTCSANOW:
		cmd = unix.TIOCSETA
	case tioTCSADRAIN:
		cmd = unix.TIOCSETAW
	case tioTCSAFLUSH:
		cmd = unix.TIOCSETAF
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
	var what int
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(unix.TIOCFLUSH),
		uintptr(unsafe.Pointer(&what)),
	)
	if errno == 0 {
		return nil
	}
	return fmt.Errorf("tcflush: %v", errno)
}
