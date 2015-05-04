package util

import (
	"os"
)

func ReadByte() (byte, error) {
	b := make([]byte, 1)

	if _, err := os.Stdin.Read(b); err != nil {
		return 0, err
	} else {
		return b[0], nil
	}
}

func WriteByte(b byte) error {
	if _, err := os.Stdout.Write([]byte{b}); err != nil {
		return err
	}
	return nil
}
