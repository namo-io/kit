package key

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	RequestId      = Key("request.id")
	Authorization  = Key("authorization")
	ServiceName    = Key("service.name")
	ServiceId      = Key("service.id")
	ServiceVersion = Key("service.version")
)

var Keys = []Key{
	RequestId,
	Authorization,
	ServiceName,
	ServiceId,
	ServiceVersion,
}
