package exception

type (
	HttpResponseError interface {
		StatusCode() int
		Error() string
	}

	BadRequest struct {
		Message string
	}
	NotFound struct {
		Message string
	}
)

func (error BadRequest) Error() string {
	return error.Message
}

func (error BadRequest) StatusCode() int {
	return 400
}

func (error NotFound) Error() string {
	return error.Message
}

func (error NotFound) StatusCode() int {
	return 404
}
