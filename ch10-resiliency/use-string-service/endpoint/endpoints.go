package endpoint

type UseStringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type UseStringResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}
