package onet

import (
	"errors"
	"fmt"
	"strings"
)

/*
   @Author: orbit-w
   @File: error
   @2024 4月 周五 21:12
*/

func IsClosedConnError(err error) bool {
	/*
		`use of closed file or network connection` (Go ver > 1.8, internal/pool.ErrClosing)
		`mux: listener closed` (cmux.ErrListenerClosed)
	*/
	return err != nil && strings.Contains(err.Error(), "closed")
}

func ExceedMaxIncomingPacket(size uint32) error {
	return errors.New(fmt.Sprintf("exceed max incoming packet size: %d", size))
}

func ReadBodyFailed(err error) error {
	return errors.New(fmt.Sprintf("read body failed: %s", err.Error()))
}
