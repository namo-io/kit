package mctx

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	RequestIdKey     = Key("request.id")
	AuthorizationKey = Key("authorization")
	AppNameKey       = Key("app.name")
	AppIdKey         = Key("app.id")
	AppVersionKey    = Key("app.version")
)

var Keys = []Key{
	RequestIdKey,
	AuthorizationKey,
	AppNameKey,
	AppIdKey,
	AppVersionKey,
}
