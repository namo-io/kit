package log

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/namo-io/kit/pkg/log/logger/typist"
	hooks "github.com/namo-io/kit/pkg/log/logger/typist/hookers"
	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	Debug("TEST")
	Error("TEST")
	q := WithField("QWE", 123).WithField("T", 1)
	q.Error("TEST")
	WithField("QWE", 123).Errorf("TEST %v", "QWE")
}

func TestHookerInner(t *testing.T) {
	typist := typist.New()
	typist.AddHooker(hooks.NewTestHooker(false))
	eshooker, err := hooks.NewElasticSearchSender(&hooks.ElasticSearchSenderConfig{
		Endpoints:       []string{"http://localhost:9200"},
		Sniff:           false,
		Index:           "log-test",
		IsCallStackFire: true,
	})
	if err != nil {
		typist.Error(err)
		return
	}
	typist.AddHooker(eshooker)

	typist.WithField("Field1", "Hello").WithField("app", "spc-cicd-api").Debug("TEST")
	count := 0
	for {
		count++
		time.Sleep(time.Second * 2)
		typist.WithField("Field1", "Hello").WithField("app", "spc-cicd-api").Errorf("TEST, %v", count)
	}
}

func TestHooker(t *testing.T) {
	TestHookerInner(t)
}

func TestHookerError(t *testing.T) {
	typist := typist.New()
	typist.AddHooker(hooks.NewTestHooker(true))
	typist.Debug("TEST")
	typist.Debug("TEST")
}

func TestBenchMark(t *testing.T) {

	var l sync.WaitGroup

	prev := time.Now()
	for i := 0; i < 10; i++ {
		me := fmt.Sprintf("routine-%v", i)
		l.Add(1)

		go func() {
			_count := 0
			for _count < 10000 {
				_count++
				WithField("TEST", 1).Infof("%v: TEST-%v", me, _count)
			}

			l.Done()
		}()
	}

	l.Wait()

	now := time.Now()

	logrus := logrus.New()
	prevLogrus := time.Now()
	for i := 0; i < 10; i++ {
		me := fmt.Sprintf("routine-%v", i)
		l.Add(1)

		go func() {
			_count := 0
			for _count < 10000 {
				_count++
				logrus.WithField("TEST", 1).Infof("%v: TEST-%v", me, _count)
			}

			l.Done()
		}()
	}

	l.Wait()

	nowLogrus := time.Now()
	fmt.Printf("typist: %v - %v = %v\r\n", prev.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Sub(prev).Milliseconds())
	fmt.Printf("logrus: %v - %v = %v\r\n", prevLogrus.Format(time.RFC3339Nano), nowLogrus.Format(time.RFC3339Nano), nowLogrus.Sub(prevLogrus).Milliseconds())
}
