package models

type URLRequest struct {
	Url string `json:"url"`
}

type URLResponse struct {
	Result string `json:"result"`
}