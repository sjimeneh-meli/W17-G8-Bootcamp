package seeders

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"log"
)

type seeder struct {
	productService services.ProductServiceI
}

func NewSeeder(productService services.ProductServiceI) *seeder {
	return &seeder{
		productService: productService,
	}
}

func (s *seeder) LoadAllData() {
	s.LoadProducts()
}

func (s *seeder) LoadProducts() {

	products := []models.Product{
		{
			Id:                             1,
			ProductCode:                    "APP-GRN-001",
			Description:                    "Manzana Verde Granny Smith",
			Width:                          8.0,
			Height:                         8.0,
			Length:                         8.0,
			NetWeight:                      0.18,
			ExpirationRate:                 0.05, // 5% por día
			RecommendedFreezingTemperature: -1.0, // No recomendado congelar, pero si se hiciera a -1°C
			FreezingRate:                   0.0,  // No aplica, producto fresco
			ProductTypeID:                  1,    // Fruta
			SellerID:                       nil,
		},
		{
			Id:                             2,
			ProductCode:                    "BAN-YEL-001",
			Description:                    "Banano Cavendish",
			Width:                          4.0,
			Height:                         20.0,
			Length:                         4.0,
			NetWeight:                      0.15,
			ExpirationRate:                 0.1, // 10% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruta
			SellerID:                       nil,
		},
		{
			Id:                             3,
			ProductCode:                    "TOM-RED-001",
			Description:                    "Tomate Chonto",
			Width:                          7.0,
			Height:                         7.0,
			Length:                         7.0,
			NetWeight:                      0.12,
			ExpirationRate:                 0.08, // 8% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura
			SellerID:                       nil,
		},
		{
			Id:                             4,
			ProductCode:                    "CAR-ORG-001",
			Description:                    "Zanahoria Fresca",
			Width:                          3.0,
			Height:                         18.0,
			Length:                         3.0,
			NetWeight:                      0.09,
			ExpirationRate:                 0.03, // 3% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura
			SellerID:                       nil,
		},
		{
			Id:                             5,
			ProductCode:                    "MIL-WHO-001",
			Description:                    "Leche Entera Pasteurizada 1L",
			Width:                          7.0,
			Height:                         25.0,
			Length:                         7.0,
			NetWeight:                      1.03,
			ExpirationRate:                 0.15, // 15% por día
			RecommendedFreezingTemperature: -0.5,
			FreezingRate:                   0.0,
			ProductTypeID:                  3, // Lácteo
			SellerID:                       nil,
		},
		{
			Id:                             6,
			ProductCode:                    "EGG-WHT-001",
			Description:                    "Cartón de Huevos Blancos (12 unidades)",
			Width:                          15.0,
			Height:                         7.0,
			Length:                         25.0,
			NetWeight:                      0.75,
			ExpirationRate:                 0.02, // 2% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  4, // Proteína Animal (Huevos)
			SellerID:                       nil,
		},
		{
			Id:                             7,
			ProductCode:                    "CHI-BRE-001",
			Description:                    "Pechuga de Pollo Fresca (500g)",
			Width:                          15.0,
			Height:                         5.0,
			Length:                         20.0,
			NetWeight:                      0.5,
			ExpirationRate:                 0.3, // 30% por día
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.5, // Tasa de congelación rápida
			ProductTypeID:                  5,   // Carne
			SellerID:                       nil,
		},
		{
			Id:                             8,
			ProductCode:                    "BEE-STE-001",
			Description:                    "Carne de Res para Asar (Lomo) (500g)",
			Width:                          15.0,
			Height:                         5.0,
			Length:                         25.0,
			NetWeight:                      0.5,
			ExpirationRate:                 0.25, // 25% por día
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.4,
			ProductTypeID:                  5, // Carne
			SellerID:                       nil,
		},
		{
			Id:                             9,
			ProductCode:                    "LET-GRN-001",
			Description:                    "Lechuga Romana",
			Width:                          20.0,
			Height:                         15.0,
			Length:                         20.0,
			NetWeight:                      0.3,
			ExpirationRate:                 0.12, // 12% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura
			SellerID:                       nil,
		},
		{
			Id:                             10,
			ProductCode:                    "ORA-JUI-001",
			Description:                    "Naranja para Jugo",
			Width:                          9.0,
			Height:                         9.0,
			Length:                         9.0,
			NetWeight:                      0.22,
			ExpirationRate:                 0.07, // 7% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruta
			SellerID:                       nil,
		},
		{
			Id:                             11,
			ProductCode:                    "POT-BRO-001",
			Description:                    "Papa Pastusa",
			Width:                          8.0,
			Height:                         8.0,
			Length:                         8.0,
			NetWeight:                      0.25,
			ExpirationRate:                 0.01, // 1% por día (larga duración)
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura (Tubérculo)
			SellerID:                       nil,
		},
		{
			Id:                             12,
			ProductCode:                    "BRO-GRN-001",
			Description:                    "Brócoli Fresco",
			Width:                          15.0,
			Height:                         15.0,
			Length:                         15.0,
			NetWeight:                      0.4,
			ExpirationRate:                 0.09, // 9% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura
			SellerID:                       nil,
		},
		{
			Id:                             13,
			ProductCode:                    "YOG-NAT-001",
			Description:                    "Yogurt Natural (500g)",
			Width:                          10.0,
			Height:                         12.0,
			Length:                         10.0,
			NetWeight:                      0.5,
			ExpirationRate:                 0.18, // 18% por día
			RecommendedFreezingTemperature: -0.5,
			FreezingRate:                   0.0,
			ProductTypeID:                  3, // Lácteo
			SellerID:                       nil,
		},
		{
			Id:                             14,
			ProductCode:                    "CHE-MOS-001",
			Description:                    "Queso Mozzarella (250g)",
			Width:                          10.0,
			Height:                         3.0,
			Length:                         15.0,
			NetWeight:                      0.25,
			ExpirationRate:                 0.04, // 4% por día
			RecommendedFreezingTemperature: -2.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  3, // Lácteo
			SellerID:                       nil,
		},
		{
			Id:                             15,
			ProductCode:                    "SAL-FRE-001",
			Description:                    "Salmón Fresco Filete (200g)",
			Width:                          10.0,
			Height:                         3.0,
			Length:                         20.0,
			NetWeight:                      0.2,
			ExpirationRate:                 0.4, // 40% por día (muy perecedero)
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.6,
			ProductTypeID:                  6, // Pescado
			SellerID:                       nil,
		},
		{
			Id:                             16,
			ProductCode:                    "PIG-CHO-001",
			Description:                    "Chuleta de Cerdo Fresca (250g)",
			Width:                          12.0,
			Height:                         3.0,
			Length:                         18.0,
			NetWeight:                      0.25,
			ExpirationRate:                 0.35, // 35% por día
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.45,
			ProductTypeID:                  5, // Carne
			SellerID:                       nil,
		},
		{
			Id:                             17,
			ProductCode:                    "CUC-GRN-001",
			Description:                    "Pepino Cohombro",
			Width:                          5.0,
			Height:                         25.0,
			Length:                         5.0,
			NetWeight:                      0.2,
			ExpirationRate:                 0.06, // 6% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura
			SellerID:                       nil,
		},
		{
			Id:                             18,
			ProductCode:                    "BER-MIX-001",
			Description:                    "Mix de Berries Frescos (200g)",
			Width:                          12.0,
			Height:                         5.0,
			Length:                         12.0,
			NetWeight:                      0.2,
			ExpirationRate:                 0.2, // 20% por día (muy delicado)
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruta
			SellerID:                       nil,
		},
		{
			Id:                             19,
			ProductCode:                    "GAR-WHT-001",
			Description:                    "Cabeza de Ajo",
			Width:                          6.0,
			Height:                         5.0,
			Length:                         6.0,
			NetWeight:                      0.05,
			ExpirationRate:                 0.005, // 0.5% por día (muy duradero)
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Verdura (Condimento)
			SellerID:                       nil,
		},
		{
			Id:                             20,
			ProductCode:                    "STR-RED-001",
			Description:                    "Fresa Fresca (250g)",
			Width:                          15.0,
			Height:                         5.0,
			Length:                         15.0,
			NetWeight:                      0.25,
			ExpirationRate:                 0.25, // 25% por día
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruta
			SellerID:                       nil,
		},
	}

	currentProducts, _ := s.productService.GetAll()
	if len(currentProducts) == 0 {
		_, err := s.productService.CreateByBatch(products)
		if err != nil {
			log.Println(fmt.Sprintf(err.Error()))
		}
	}
}
