package web

type Response struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}