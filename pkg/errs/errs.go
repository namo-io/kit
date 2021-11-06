package errs

type Error struct {
	// error type ex) IAM/CreateUser
	Domain string

	// error idenfitifer ex) UID_PROPERTY_IS_STRANGE
	Reason string

	// error message ex) UID property is strange
	Message string

	Extensions map[string]string
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) WithExt(key string, value string) Error {
	if e.Extensions == nil {
		e.Extensions = map[string]string{}
	}

	e.Extensions[key] = value
	return e
}

func New(domain string, reason string, message string) Error {
	return Error{
		Domain:     domain,
		Reason:     reason,
		Message:    message,
		Extensions: map[string]string{},
	}
}

func WithErr(domain string, reason string, err error) Error {
	return New(domain, reason, err.Error())
}
