package requests

type BuyerRequest struct {
	CardNumberId string `json:"id_card_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}
