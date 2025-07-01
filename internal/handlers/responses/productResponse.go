package responses

type ProductResponse struct {
	Description                    string  `json:"description"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	FreezingRate                   float64 `json:"freezing_rate"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ProductCode                    string  `json:"product_code"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	Width                          float64 `json:"width"`
	ProductTypeID                  int     `json:"product_type_id"`
	SellerID                       *int    `json:"seller_id,omitempty"` // Usamos un puntero para indicar que no es obligatorio
}
