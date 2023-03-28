package util

type ApiResponse struct {
	Error ErrorDescription `json:"error"`
}

type ErrorDescription struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
