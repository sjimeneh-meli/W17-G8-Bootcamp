package responses

type LocalitySellerReport struct {
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	SellerCount  int    `json:"sellers_count"`
}

type LocalitySellerReportResponse struct {
	Data []LocalitySellerReport `json:"data"`
}
