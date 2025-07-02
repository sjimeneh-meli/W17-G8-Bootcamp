package responses

type DataResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
