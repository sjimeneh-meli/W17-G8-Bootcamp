package responses

type BuyerResponse struct {
	Id           int    `json:"id"`
	CardNumberId string `json:"id_card_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}
