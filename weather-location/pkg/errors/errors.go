package errors

type Error struct {
	StatusCode int
	Message    string
}

func New(message string, statusCode int) Error {
	return Error{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e Error) Error() string {
	return e.Message
}
