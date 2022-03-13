package log

import "testing"

func TestLog(t *testing.T) {
	Info("TEST")
	Error("TEST")
	Warn("TEST")
	Debug("TEST")
	Trace("TEST")
	Trace("TEST")
	Trace("TEST")
	Trace("TEST")
	test(gLog)
	gLog.WithField("app.name", "Test").Info("QWE")
	gLog.WithField("addr", "127.0.0.1").Debug("QQ")
	gLog.WithFields(map[string]string{
		"addr":  "127.0.0.1",
		"addr2": "127.0.0.2",
	}).Debug("QQ")
}

func test(l Log) {
	l.Debug("Q")
}
