package httpapi

type healthResponse struct {
	Status string `json:"status"`
}

type countsResponse struct {
	Views int64 `json:"views"`
	Likes int64 `json:"likes"`
}
