package log

import "testing"

func TestLog(t *testing.T) {
	Debug("TEST")
	Info("TEST")
	test(glog)
	glog.WithField("addr", "127.0.0.1").Debug("QQ")
}

func test(l Log) {
	l.Debug("Q")
}
