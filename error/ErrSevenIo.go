package error

import "fmt"

type ErrSevenIo string

// Error implements error.
func (e ErrSevenIo) Error() string {
	return fmt.Sprintf("ErrSevenIo: %s", string(e))
}

func NewErrSevenIo(format string, args ...any) error {
	return ErrSevenIo(fmt.Sprintf(format, args...))
}
