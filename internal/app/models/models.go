package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	URL string `json:"result"`
}
