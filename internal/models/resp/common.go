package resp

type ResponseError struct {
	Message string `json:"message"`
}

type SimpleResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
