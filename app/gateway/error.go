package gateway

import (
	"errors"
)

/*
   @Author: orbit-w
   @File: error
   @2023 12月 周六 22:27
*/

func AgentDecodePatternErr(pattern string) error {
	return errors.New("agent decode pattern error: " + pattern)
}
