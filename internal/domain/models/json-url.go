package models

// URLRequest URL для сокращения в формате JSON.
type URLRequest struct {
	URL string `json:"url"`
}

// URLResponse ответ сокращенного URL в формате JSON.
type URLResponse struct {
	Result string `json:"result"`
}
