package requests

type ProductRequest struct {
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Width                          float64 `json:"width"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductTypeID                  int64   `json:"product_type_id"`
	SellerID                       *int64  `json:"seller_id,omitempty"` // Usamos un puntero para indicar que no es obligatorio
}
