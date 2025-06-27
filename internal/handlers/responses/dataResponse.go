package responses

type DataResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   bool   `json:"error"`
}

func (r *DataResponse) SetError(msgErr string) {
	r.Message = msgErr
	r.Error = true
}
