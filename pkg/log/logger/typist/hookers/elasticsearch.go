package hooks

import (
	"context"
	"fmt"
	"time"

	"github.com/barkimedes/go-deepcopy"
	"github.com/namo-io/kit/pkg/log/logger/typist"
	"github.com/olivere/elastic/v7"
)

type elasticSearchSender struct {
	client          *elastic.Client
	index           string
	isCallStackFire bool
}

type ElasticSearchSenderConfig struct {
	Endpoints       []string
	Sniff           bool
	Index           string
	IsCallStackFire bool
}

type elasticSearchSenderMessage struct {
	Timestamp string                 `json:"@timestamp"`
	Level     string                 `json:"level,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

var (
	// ErrCannotCreateIndex Fired if the index is not created
	ErrCannotCreateIndex = fmt.Errorf("cannot create index")
)

const (
	CallStackKey = "Callstack"
)

func NewElasticSearchSender(cfg *ElasticSearchSenderConfig) (typist.Hooker, error) {
	// create es client
	client, err := elastic.NewClient(elastic.SetURL(cfg.Endpoints...), elastic.SetSniff(cfg.Sniff))
	if err != nil {
		return nil, err
	}

	// find & create es index
	ctx := context.Background()
	exists, err := client.IndexExists(cfg.Index).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		createIndex, err := client.CreateIndex(cfg.Index).Do(ctx)
		if err != nil {
			return nil, err
		}
		if !createIndex.Acknowledged {
			return nil, ErrCannotCreateIndex
		}
	}

	return &elasticSearchSender{
		client:          client,
		index:           cfg.Index,
		isCallStackFire: cfg.IsCallStackFire,
	}, nil
}

func (t *elasticSearchSender) Name() string {
	return "ElasticSearchSender"
}

func (e *elasticSearchSender) Fire(ctx context.Context, level typist.Level, rs *typist.Record) error {
	if rs.Level >= typist.TraceLevel {
		return nil
	}

	go e.syncFire(ctx, rs)
	return nil
}

func (e *elasticSearchSender) syncFire(ctx context.Context, rs *typist.Record) {
	msg := e.newElasticSearchSenderMessage(rs)
	_, err := e.client.Index().Index(e.index).BodyJson(msg).Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
}

func (e *elasticSearchSender) newElasticSearchSenderMessage(rs *typist.Record) *elasticSearchSenderMessage {
	meta := deepcopy.MustAnything(rs.Meta).(map[string]interface{})

	if e.isCallStackFire && rs.Level <= typist.ErrorLevel {
		str := ""
		if len(rs.Frames) > 0 {
			for _, frame := range rs.Frames {
				if frame.Func == nil {
					continue
				}
				str += fmt.Sprintf("%v (%v)\r\n", frame.Func.Name(), frame.Func.Entry())
			}
		}

		meta[CallStackKey] = str
	}

	return &elasticSearchSenderMessage{
		Timestamp: rs.Time.UTC().Format(time.RFC3339Nano),
		Message:   rs.Message,
		Level:     rs.Level.String(),
		Meta:      meta,
	}
}
