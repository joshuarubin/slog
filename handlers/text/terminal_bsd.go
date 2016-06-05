// +build darwin freebsd openbsd netbsd dragonfly

package text

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios
