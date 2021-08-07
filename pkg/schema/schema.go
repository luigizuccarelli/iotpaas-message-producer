package schema

// Response schema
type Response struct {
	StatusCode string `json:"statuscode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type IOTPaaS struct {
	Id string `json:"Id"`
}
