package log

import "testing"

func TestLog(t *testing.T) {
	Debug("TEST")
	Info("TEST")
	test(glog)
}

func test(l Log) {
	l.Debug("Q")
}
