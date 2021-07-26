package keys

type Key string

const (
	RequestID = Key("RequestID")
)

func (k Key) String() string {
	return string(k)
}
