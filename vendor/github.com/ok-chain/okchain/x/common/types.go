package common

type ListResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"perPage"`
	Data    interface{} `json:"data"`
}

func GetListResponse(total, page, perPage int, data interface{}) *ListResponse {
	return &ListResponse{0, "", total, page, perPage, data}
}
