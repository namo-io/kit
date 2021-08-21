package keys

const (
	RequestId     = "request-id"
	Authorization = "authorization"
)

func All() []string {
	return []string{
		RequestId,
		Authorization,
	}
}
