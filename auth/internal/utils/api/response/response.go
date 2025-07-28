package responses

type Responce struct {
	Status string `json:"status"`
	Error  string `json:"errors,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func Ok() Responce {
	return Responce{
		Status: StatusOK,
	}
}

func Error(msg string) Responce {
	return Responce{
		Status: StatusError,
		Error:  msg,
	}
}
