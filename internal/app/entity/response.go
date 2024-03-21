package entity

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status string
	Error  error
}

type URLResponse struct {
	Response
	URL URL
}

func OKResponse() Response {
	return Response{
		Status: StatusOK,
	}
}

func ErrorResponse(err error) Response {
	return Response{
		Status: StatusError,
		Error: err,
	}
}

func OKURLResponse(url URL) URLResponse {
	return URLResponse{
		Response: OKResponse(),
		URL: url,
	}
}

func ErrorURLResponse(err error) URLResponse {
	return URLResponse{
		Response: ErrorResponse(err),
		URL: URL{},
	}
}