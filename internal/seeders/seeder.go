package seeders

import (
	"log"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

type SeederI interface {
	LoadAllData()
}

// seeder is responsible for loading initial data into the database
// seeder se encarga de cargar datos iniciales en la base de datos
type seeder struct {
	productService services.ProductServiceI
}

// NewSeeder creates a new seeder instance with the provided product service
// NewSeeder crea una nueva instancia del seeder con el servicio de productos proporcionado
func NewSeeder(productService services.ProductServiceI) SeederI {
	return &seeder{
		productService: productService,
	}
}

// LoadAllData orchestrates the loading of all seed data
// This method serves as the main entry point for data seeding operations
// LoadAllData orquesta la carga de todos los datos semilla
// Este método sirve como punto de entrada principal para las operaciones de siembra de datos
func (s *seeder) LoadAllData() {
	s.LoadProducts()
}

// LoadProducts loads initial product catalog into the database
// Only executes if no products exist to prevent data duplication
// LoadProducts carga el catálogo inicial de productos en la base de datos
// Solo se ejecuta si no existen productos para prevenir duplicación de datos
func (s *seeder) LoadProducts() {

	// Seed data: comprehensive product catalog covering various food categories
	// Datos semilla: catálogo integral de productos cubriendo varias categorías de alimentos
	products := []models.Product{
		{
			// Fresh Fruit Category - Categoría de Frutas Frescas
			Id:                             1,
			ProductCode:                    "APP-GRN-001",
			Description:                    "Manzana Verde Granny Smith",
			Width:                          8.0,
			Height:                         8.0,
			Length:                         8.0,
			NetWeight:                      0.18,
			ExpirationRate:                 0.05, // 5% daily deterioration rate - Tasa de deterioro diario del 5%
			RecommendedFreezingTemperature: -1.0, // Not recommended for freezing - No recomendado para congelar
			FreezingRate:                   0.0,  // Fresh product, no freezing rate - Producto fresco, sin tasa de congelación
			ProductTypeID:                  1,    // Fruit category - Categoría fruta
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
			ExpirationRate:                 0.1, // High perishability - Alta perecibilidad
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruit category - Categoría fruta
			SellerID:                       nil,
		},
		{
			// Vegetable Category - Categoría de Verduras
			Id:                             3,
			ProductCode:                    "TOM-RED-001",
			Description:                    "Tomate Chonto",
			Width:                          7.0,
			Height:                         7.0,
			Length:                         7.0,
			NetWeight:                      0.12,
			ExpirationRate:                 0.08, // Moderate perishability - Perecibilidad moderada
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable category - Categoría verdura
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
			ExpirationRate:                 0.03, // Low deterioration rate - Baja tasa de deterioro
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable category - Categoría verdura
			SellerID:                       nil,
		},
		{
			// Dairy Products Category - Categoría de Productos Lácteos
			Id:                             5,
			ProductCode:                    "MIL-WHO-001",
			Description:                    "Leche Entera Pasteurizada 1L",
			Width:                          7.0,
			Height:                         25.0,
			Length:                         7.0,
			NetWeight:                      1.03,
			ExpirationRate:                 0.15, // High perishability due to dairy nature - Alta perecibilidad por naturaleza láctea
			RecommendedFreezingTemperature: -0.5,
			FreezingRate:                   0.0,
			ProductTypeID:                  3, // Dairy category - Categoría lácteo
			SellerID:                       nil,
		},
		{
			// Animal Protein Category - Categoría de Proteína Animal
			Id:                             6,
			ProductCode:                    "EGG-WHT-001",
			Description:                    "Cartón de Huevos Blancos (12 unidades)",
			Width:                          15.0,
			Height:                         7.0,
			Length:                         25.0,
			NetWeight:                      0.75,
			ExpirationRate:                 0.02, // Low deterioration due to shell protection - Bajo deterioro por protección del cascarón
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  4, // Animal protein (Eggs) - Proteína animal (Huevos)
			SellerID:                       nil,
		},
		{
			// Meat Category - Categoría de Carnes
			Id:                             7,
			ProductCode:                    "CHI-BRE-001",
			Description:                    "Pechuga de Pollo Fresca (500g)",
			Width:                          15.0,
			Height:                         5.0,
			Length:                         20.0,
			NetWeight:                      0.5,
			ExpirationRate:                 0.3, // Very high perishability - Muy alta perecibilidad
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.5, // Fast freezing rate for meat preservation - Tasa de congelación rápida para preservación
			ProductTypeID:                  5,   // Meat category - Categoría carne
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
			ExpirationRate:                 0.25, // High perishability typical of red meat - Alta perecibilidad típica de carne roja
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.4,
			ProductTypeID:                  5, // Meat category - Categoría carne
			SellerID:                       nil,
		},
		{
			// Additional Vegetables - Verduras Adicionales
			Id:                             9,
			ProductCode:                    "LET-GRN-001",
			Description:                    "Lechuga Romana",
			Width:                          20.0,
			Height:                         15.0,
			Length:                         20.0,
			NetWeight:                      0.3,
			ExpirationRate:                 0.12, // Leafy greens have moderate deterioration - Hojas verdes tienen deterioro moderado
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable category - Categoría verdura
			SellerID:                       nil,
		},
		{
			// Citrus Fruits - Frutas Cítricas
			Id:                             10,
			ProductCode:                    "ORA-JUI-001",
			Description:                    "Naranja para Jugo",
			Width:                          9.0,
			Height:                         9.0,
			Length:                         9.0,
			NetWeight:                      0.22,
			ExpirationRate:                 0.07, // Citrus fruits have moderate shelf life - Frutas cítricas tienen vida útil moderada
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruit category - Categoría fruta
			SellerID:                       nil,
		},
		{
			// Root Vegetables - Verduras de Raíz
			Id:                             11,
			ProductCode:                    "POT-BRO-001",
			Description:                    "Papa Pastusa",
			Width:                          8.0,
			Height:                         8.0,
			Length:                         8.0,
			NetWeight:                      0.25,
			ExpirationRate:                 0.01, // Very low deterioration, long shelf life - Muy bajo deterioro, larga vida útil
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable (Tuber) - Verdura (Tubérculo)
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
			ExpirationRate:                 0.09, // Cruciferous vegetables have moderate perishability - Verduras crucíferas tienen perecibilidad moderada
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable category - Categoría verdura
			SellerID:                       nil,
		},
		{
			// Additional Dairy Products - Productos Lácteos Adicionales
			Id:                             13,
			ProductCode:                    "YOG-NAT-001",
			Description:                    "Yogurt Natural (500g)",
			Width:                          10.0,
			Height:                         12.0,
			Length:                         10.0,
			NetWeight:                      0.5,
			ExpirationRate:                 0.18, // Fermented dairy has higher perishability - Lácteos fermentados tienen mayor perecibilidad
			RecommendedFreezingTemperature: -0.5,
			FreezingRate:                   0.0,
			ProductTypeID:                  3, // Dairy category - Categoría lácteo
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
			ExpirationRate:                 0.04, // Lower moisture cheese has better shelf life - Queso con menor humedad tiene mejor vida útil
			RecommendedFreezingTemperature: -2.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  3, // Dairy category - Categoría lácteo
			SellerID:                       nil,
		},
		{
			// Seafood Category - Categoría de Mariscos
			Id:                             15,
			ProductCode:                    "SAL-FRE-001",
			Description:                    "Salmón Fresco Filete (200g)",
			Width:                          10.0,
			Height:                         3.0,
			Length:                         20.0,
			NetWeight:                      0.2,
			ExpirationRate:                 0.4, // Fish has extremely high perishability - El pescado tiene perecibilidad extremadamente alta
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.6,
			ProductTypeID:                  6, // Fish category - Categoría pescado
			SellerID:                       nil,
		},
		{
			// Pork Products - Productos de Cerdo
			Id:                             16,
			ProductCode:                    "PIG-CHO-001",
			Description:                    "Chuleta de Cerdo Fresca (250g)",
			Width:                          12.0,
			Height:                         3.0,
			Length:                         18.0,
			NetWeight:                      0.25,
			ExpirationRate:                 0.35, // Pork has high perishability similar to other meats - Cerdo tiene alta perecibilidad similar a otras carnes
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.45,
			ProductTypeID:                  5, // Meat category - Categoría carne
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
			ExpirationRate:                 0.06, // Cucumbers have moderate deterioration - Pepinos tienen deterioro moderado
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable category - Categoría verdura
			SellerID:                       nil,
		},
		{
			// Delicate Fruits - Frutas Delicadas
			Id:                             18,
			ProductCode:                    "BER-MIX-001",
			Description:                    "Mix de Berries Frescos (200g)",
			Width:                          12.0,
			Height:                         5.0,
			Length:                         12.0,
			NetWeight:                      0.2,
			ExpirationRate:                 0.2, // Berries are extremely delicate and perishable - Las berries son extremadamente delicadas y perecederas
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruit category - Categoría fruta
			SellerID:                       nil,
		},
		{
			// Aromatics and Seasonings - Aromáticos y Condimentos
			Id:                             19,
			ProductCode:                    "GAR-WHT-001",
			Description:                    "Cabeza de Ajo",
			Width:                          6.0,
			Height:                         5.0,
			Length:                         6.0,
			NetWeight:                      0.05,
			ExpirationRate:                 0.005, // Garlic has excellent shelf life - El ajo tiene excelente vida útil
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  2, // Vegetable (Seasoning) - Verdura (Condimento)
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
			ExpirationRate:                 0.25, // Strawberries are highly perishable - Las fresas son altamente perecederas
			RecommendedFreezingTemperature: -1.0,
			FreezingRate:                   0.0,
			ProductTypeID:                  1, // Fruit category - Categoría fruta
			SellerID:                       nil,
		},
	}

	// Check if products already exist to prevent duplication
	// Verificar si ya existen productos para prevenir duplicación
	currentProducts, _ := s.productService.GetAll()
	if len(currentProducts) == 0 {
		// Batch create all products for better performance
		// Crear todos los productos en lote para mejor rendimiento
		_, err := s.productService.CreateByBatch(products)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
