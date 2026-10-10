package entities

type DataResponse struct {
	Code       int                    `json:"code"`
	Data       interface{}            `json:"data"`
	Message    string                 `json:"message"`
	DataSchema map[string]interface{} `json:"data_schema"`
}

type DataListResponse struct {
	Code       int                    `json:"code"`
	Data       interface{}            `json:"data"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPage  int                    `json:"total_page"`
	TotalData  int64                  `json:"total_data"`
	Message    string                 `json:"message"`
	DataSchema map[string]interface{} `json:"data_schema"`
}
