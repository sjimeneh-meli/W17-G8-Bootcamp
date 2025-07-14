-- =================================================================
-- SCRIPT DE CREACIÓN DE LA BASE DE DATOS Y TABLAS
-- =================================================================

DROP DATABASE IF EXISTS productos_frescos;

CREATE DATABASE productos_fescos;

USE productos_fescos;

-- Eliminación de tablas en orden inverso para evitar conflictos de claves foráneas
DROP TABLE IF EXISTS `user_rol`;
DROP TABLE IF EXISTS `purchase_orders`;
DROP TABLE IF EXISTS `inbound_orders`;
DROP TABLE IF EXISTS `product_batches`;
DROP TABLE IF EXISTS `product_records`;
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `sections`;
DROP TABLE IF EXISTS `products_types`;
DROP TABLE IF EXISTS `employees`;
DROP TABLE IF EXISTS `warehouse`;
DROP TABLE IF EXISTS `carriers`;
DROP TABLE IF EXISTS `sellers`;
DROP TABLE IF EXISTS `localities`;
DROP TABLE IF EXISTS `provinces`;
DROP TABLE IF EXISTS `countries`;
DROP TABLE IF EXISTS `buyers`;
DROP TABLE IF EXISTS `order_status`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `rol`;

-- Creación de la tabla 'countries'
CREATE TABLE `countries` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `country_name` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'provinces'
CREATE TABLE `provinces` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `province_name` VARCHAR(255) NOT NULL,
  `id_country_fk` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id_country_fk`) REFERENCES `countries`(`id`)
);

-- Creación de la tabla 'localities'
CREATE TABLE `localities` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `locality_name` VARCHAR(255) NOT NULL,
  `province_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`province_id`) REFERENCES `provinces`(`id`)
);

-- Creación de la tabla 'sellers'
CREATE TABLE `sellers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `cid` VARCHAR(255) NOT NULL,
  `company_name` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `telephone` VARCHAR(255) NOT NULL,
  `locality_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`locality_id`) REFERENCES `localities`(`id`)
);

-- Creación de la tabla 'buyers'
CREATE TABLE `buyers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `id_card_number` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(255) NOT NULL,
  `last_name` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'warehouse'
CREATE TABLE `warehouse` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `address` VARCHAR(255) NOT NULL,
  `telephone` VARCHAR(255) NOT NULL,
  `warehouse_code` VARCHAR(255) NOT NULL,
  `locality_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`locality_id`) REFERENCES `localities`(`id`)
);

-- Creación de la tabla 'employees'
CREATE TABLE `employees` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `id_card_number` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(255) NOT NULL,
  `last_name` VARCHAR(255) NOT NULL,
  `warehouse_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`)
);

-- Creación de la tabla 'products_types'
CREATE TABLE `products_types` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'sections'
CREATE TABLE `sections` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `section_number` VARCHAR(255) NOT NULL,
  `current_capacity` INT,
  `current_temperature` DECIMAL(19,2),
  `maximum_capacity` INT,
  `minimum_capacity` INT,
  `minimum_temperature` DECIMAL(19,2),
  `warehouse_id` INT NOT NULL,
  `product_type_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`),
  FOREIGN KEY (`product_type_id`) REFERENCES `products_types`(`id`)
);

-- Creación de la tabla 'products'
CREATE TABLE `products` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(255),
  `expiration_rate` DECIMAL(19,2),
  `freezing_rate` DECIMAL(19,2),
  `height` DECIMAL(19,2),
  `length` DECIMAL(19,2),
  `net_weight` DECIMAL(19,2),
  `product_code` VARCHAR(255) NOT NULL,
  `recommended_freezing_temperature` DECIMAL(19,2),
  `width` DECIMAL(19,2),
  `product_type_id` INT,
  `seller_id` INT,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`product_type_id`) REFERENCES `products_types`(`id`),
  FOREIGN KEY (`seller_id`) REFERENCES `sellers`(`id`)
);

-- Creación de la tabla 'product_records'
CREATE TABLE `product_records` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `last_update_date` DATETIME(6),
  `purchase_price` DECIMAL(19,2),
  `sale_price` DECIMAL(19,2),
  `product_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `products`(`id`)
);

-- Creación de la tabla 'product_batches'
CREATE TABLE `product_batches` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `batch_number` VARCHAR(255) NOT NULL,
  `current_quantity` INT,
  `current_temperature` DECIMAL(19,2),
  `due_date` DATETIME(6),
  `initial_quantity` INT,
  `manufacturing_date` DATETIME(6),
  `manufacturing_hour` DATETIME(6),
  `minimum_temperature` DECIMAL(19,2),
  `product_id` INT NOT NULL,
  `section_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `products`(`id`),
  FOREIGN KEY (`section_id`) REFERENCES `sections`(`id`)
);

-- Creación de la tabla 'inbound_orders'
CREATE TABLE `inbound_orders` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_date` DATETIME(6),
  `order_number` VARCHAR(255) NOT NULL,
  `employee_id` INT NOT NULL,
  `product_batch_id` INT NOT NULL,
  `warehouse_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`employee_id`) REFERENCES `employees`(`id`),
  FOREIGN KEY (`product_batch_id`) REFERENCES `product_batches`(`id`),
  FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`)
);

-- Creación de la tabla 'carriers'
CREATE TABLE `carriers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `cid` VARCHAR(255) NOT NULL,
  `company_name` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `telephone` VARCHAR(255) NOT NULL,
  `locality_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`locality_id`) REFERENCES `localities`(`id`)
);

-- Creación de la tabla 'order_status'
CREATE TABLE `order_status` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'purchase_orders' (ESTRUCTURA ACTUALIZADA)
CREATE TABLE `purchase_orders` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_number` VARCHAR(255) NOT NULL,
  `order_date` DATETIME NOT NULL,
  `tracking_code` VARCHAR(255),
  `buyer_id` INT NOT NULL,
  `product_record_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`buyer_id`) REFERENCES `buyers`(`id`),
  FOREIGN KEY (`product_record_id`) REFERENCES `product_records`(`id`)
);

-- Creación de la tabla 'users'
CREATE TABLE `users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'rol'
CREATE TABLE `rol` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `rol_name` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'user_rol'
CREATE TABLE `user_rol` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `usuario_id` INT NOT NULL,
  `rol_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`usuario_id`) REFERENCES `users`(`id`),
  FOREIGN KEY (`rol_id`) REFERENCES `rol`(`id`)
);


-- =================================================================
-- SCRIPT DE INSERCIÓN DE DATOS
-- =================================================================

-- Insertando datos en 'countries' (3 registros solicitados)
INSERT INTO `countries` (`id`, `country_name`) VALUES
(1, 'Colombia'),
(2, 'Argentina'),
(3, 'México');

-- Insertando datos en 'provinces'
INSERT INTO `provinces` (`id`, `province_name`, `id_country_fk`) VALUES
(1, 'Cundinamarca', 1), -- Relacionado con countries.id 1 (Colombia)
(2, 'Antioquia', 1), -- Relacionado con countries.id 1 (Colombia)
(3, 'Valle del Cauca', 1), -- Relacionado con countries.id 1 (Colombia)
(4, 'Buenos Aires', 2), -- Relacionado con countries.id 2 (Argentina)
(5, 'Córdoba', 2), -- Relacionado con countries.id 2 (Argentina)
(6, 'Santa Fe', 2), -- Relacionado con countries.id 2 (Argentina)
(7, 'Jalisco', 3), -- Relacionado con countries.id 3 (México)
(8, 'Nuevo León', 3), -- Relacionado con countries.id 3 (México)
(9, 'Ciudad de México', 3); -- Relacionado con countries.id 3 (México)

-- Insertando datos en 'localities'
INSERT INTO `localities` (`id`, `locality_name`, `province_id`) VALUES
(1, 'Bogotá', 1), -- Relacionado con provinces.id 1 (Cundinamarca)
(2, 'Medellín', 2), -- Relacionado con provinces.id 2 (Antioquia)
(3, 'Cali', 3), -- Relacionado con provinces.id 3 (Valle del Cauca)
(4, 'La Plata', 4), -- Relacionado con provinces.id 4 (Buenos Aires)
(5, 'Córdoba Capital', 5), -- Relacionado con provinces.id 5 (Córdoba)
(6, 'Rosario', 6), -- Relacionado con provinces.id 6 (Santa Fe)
(7, 'Guadalajara', 7), -- Relacionado con provinces.id 7 (Jalisco)
(8, 'Monterrey', 8), -- Relacionado con provinces.id 8 (Nuevo León)
(9, 'CDMX', 9); -- Relacionado con provinces.id 9 (Ciudad de México)

-- Insertando datos en 'sellers' (5 registros solicitados)
INSERT INTO `sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES
(1, 'SEL-001', 'Distribuidora Lácteos del Campo', 'Calle Falsa 123', '555-0101', 1), -- Ubicado en localities.id 1 (Bogotá)
(2, 'SEL-002', 'Carnes de la Pampa S.A.', 'Avenida Siempreviva 742', '555-0102', 4), -- Ubicado en localities.id 4 (La Plata)
(3, 'SEL-003', 'Frutas y Verduras del Sol', 'Carrera 10 #20-30', '555-0103', 7), -- Ubicado en localities.id 7 (Guadalajara)
(4, 'SEL-004', 'Pescados del Pacífico', 'Calle del Mar 45', '555-0104', 2), -- Ubicado en localities.id 2 (Medellín)
(5, 'SEL-005', 'Congelados Express', 'Avenida de los Hielos 89', '555-0105', 8); -- Ubicado en localities.id 8 (Monterrey)

-- Insertando datos en 'buyers' (6 registros solicitados)
INSERT INTO `buyers` (`id`, `id_card_number`, `first_name`, `last_name`) VALUES
(1, '10101010', 'Ignacio', 'Garcia'),
(2, '20202020', 'Julian', 'Nahuel'),
(3, '30303030', 'Karen', 'Perez'),
(4, '40404040', 'Samuel', 'Jimenez'),
(5, '50505050', 'Gabriel', 'Lopez'),
(6, '60606060', 'Juan', 'Regino');

-- Insertando datos en 'warehouse' (10 registros solicitados)
INSERT INTO `warehouse` (`id`, `address`, `telephone`, `warehouse_code`, `locality_id`) VALUES
(1, 'Zona Franca, Bodega 10', '555-0201', 'BOG-ZF-01', 1), -- Ubicado en localities.id 1 (Bogotá)
(2, 'Parque Industrial, Nave 5', '555-0202', 'MED-PI-01', 2), -- Ubicado en localities.id 2 (Medellín)
(3, 'Central de Abastos, Bodega A2', '555-0203', 'CAL-CA-01', 3), -- Ubicado en localities.id 3 (Cali)
(4, 'Polígono Industrial Sur, Módulo 8', '555-0204', 'BUE-PIS-01', 4), -- Ubicado en localities.id 4 (La Plata)
(5, 'Centro Logístico Norte, Dock 15', '555-0205', 'COR-CLN-01', 5), -- Ubicado en localities.id 5 (Córdoba Capital)
(6, 'Puerto Norte, Bodega 3', '555-0206', 'ROS-PN-01', 6), -- Ubicado en localities.id 6 (Rosario)
(7, 'Parque Logístico Jalisco, Bodega 20', '555-0207', 'GDL-PLJ-01', 7), -- Ubicado en localities.id 7 (Guadalajara)
(8, 'Parque Industrial, Nave 8', '555-0208', 'MTY-PI-01', 8), -- Ubicado en localities.id 8 (Monterrey)
(9, 'Almacén Central, Sector 3', '555-0209', 'CDMX-AC-01', 9), -- Ubicado en localities.id 9 (CDMX)
(10, 'Bodegas del Teusaquillo', '555-0210', 'BOG-TE-02', 1); -- Ubicado en localities.id 1 (Bogotá)

-- Insertando datos en 'employees'
INSERT INTO `employees` (`id`, `id_card_number`, `first_name`, `last_name`, `warehouse_id`) VALUES
(1, '11122233A', 'Juan', 'Pérez', 1), -- Trabaja en warehouse.id 1 (Bogotá)
(2, '22233344B', 'Marta', 'García', 1), -- Trabaja en warehouse.id 1 (Bogotá)
(3, '33344455C', 'Pedro', 'Ramírez', 2), -- Trabaja en warehouse.id 2 (Medellín)
(4, '44455566D', 'Lucía', 'Fernández', 2), -- Trabaja en warehouse.id 2 (Medellín)
(5, '55566677E', 'Andrés', 'López', 3), -- Trabaja en warehouse.id 3 (Cali)
(6, '66677788F', 'Clara', 'Sanz', 3), -- Trabaja en warehouse.id 3 (Cali)
(7, '77788899G', 'Diego', 'Moreno', 4), -- Trabaja en warehouse.id 4 (La Plata)
(8, '88899900H', 'Beatriz', 'Jiménez', 4), -- Trabaja en warehouse.id 4 (La Plata)
(9, '99900011I', 'Sergio', 'Ruiz', 5), -- Trabaja en warehouse.id 5 (Córdoba)
(10, '00011122J', 'Raquel', 'Alonso', 5), -- Trabaja en warehouse.id 5 (Córdoba)
(11, '11223344K', 'Mario', 'Gutiérrez', 6), -- Trabaja en warehouse.id 6 (Rosario)
(12, '22334455L', 'Esther', 'Navarro', 6), -- Trabaja en warehouse.id 6 (Rosario)
(13, '33445566M', 'Jorge', 'Iglesias', 7), -- Trabaja en warehouse.id 7 (Guadalajara)
(14, '44556677N', 'Cristina', 'Blanco', 7), -- Trabaja en warehouse.id 7 (Guadalajara)
(15, '55667788P', 'Ricardo', 'Soto', 8), -- Trabaja en warehouse.id 8 (Monterrey)
(16, '66778899Q', 'Natalia', 'Crespo', 8), -- Trabaja en warehouse.id 8 (Monterrey)
(17, '77889900R', 'Francisco', 'Reyes', 9), -- Trabaja en warehouse.id 9 (CDMX)
(18, '88990011S', 'Verónica', 'Gil', 9), -- Trabaja en warehouse.id 9 (CDMX)
(19, '99001122T', 'Óscar', 'Ortega', 10), -- Trabaja en warehouse.id 10 (Bogotá)
(20, '00112233U', 'Mónica', 'Vega', 10); -- Trabaja en warehouse.id 10 (Bogotá)

-- Insertando datos en 'products_types'
INSERT INTO `products_types` (`id`, `description`) VALUES
(1, 'Lácteos'), (2, 'Carnes Rojas'), (3, 'Frutas'), (4, 'Pescados'), (5, 'Congelados Varios'),
(6, 'Aves'), (7, 'Panadería'), (8, 'Bebidas'), (9, 'Vinos y Licores'), (10, 'Helados'),
(11, 'Comidas Preparadas'), (12, 'Dulces y Golosinas'), (13, 'Café y Té'), (14, 'Productos Orgánicos'), (15, 'Verduras'),
(16, 'Quesos'), (17, 'Conservas'), (18, 'Snacks'), (19, 'Pastas'), (20, 'Salsas y Aderezos');

-- Insertando datos en 'sections'
INSERT INTO `sections` (`id`, `section_number`, `current_capacity`, `current_temperature`, `maximum_capacity`, `minimum_capacity`, `minimum_temperature`, `warehouse_id`, `product_type_id`) VALUES
(1, 'A-01', 100, 4.0, 200, 20, 2.0, 1, 1), -- En warehouse 1, para product_type 1 (Lácteos)
(2, 'B-01', 50, -18.0, 100, 10, -22.0, 1, 2), -- En warehouse 1, para product_type 2 (Carnes Rojas)
(3, 'A-01', 200, 10.0, 300, 30, 8.0, 2, 3), -- En warehouse 2, para product_type 3 (Frutas)
(4, 'B-01', 80, -20.0, 150, 15, -25.0, 2, 4), -- En warehouse 2, para product_type 4 (Pescados)
(5, 'A-01', 120, -18.0, 250, 25, -20.0, 3, 5), -- En warehouse 3, para product_type 5 (Congelados)
(6, 'B-01', 90, 2.0, 180, 18, 0.0, 3, 6), -- En warehouse 3, para product_type 6 (Aves)
(7, 'A-01', 150, 18.0, 200, 20, 15.0, 4, 7), -- En warehouse 4, para product_type 7 (Panadería)
(8, 'B-01', 300, 15.0, 500, 50, 12.0, 4, 8), -- En warehouse 4, para product_type 8 (Bebidas)
(9, 'A-01', 70, 14.0, 120, 12, 10.0, 5, 9), -- En warehouse 5, para product_type 9 (Vinos)
(10, 'B-01', 60, -22.0, 100, 10, -25.0, 5, 10), -- En warehouse 5, para product_type 10 (Helados)
(11, 'A-01', 85, -18.0, 150, 15, -20.0, 6, 11), -- En warehouse 6, para product_type 11 (Comidas Preparadas)
(12, 'B-01', 250, 20.0, 400, 40, 18.0, 6, 12), -- En warehouse 6, para product_type 12 (Dulces)
(13, 'A-01', 110, 22.0, 180, 18, 20.0, 7, 13), -- En warehouse 7, para product_type 13 (Café)
(14, 'B-01', 95, 12.0, 160, 16, 10.0, 7, 14), -- En warehouse 7, para product_type 14 (Orgánicos)
(15, 'A-01', 130, 8.0, 220, 22, 5.0, 8, 15), -- En warehouse 8, para product_type 15 (Verduras)
(16, 'B-01', 180, 20.0, 300, 30, 18.0, 8, 16), -- En warehouse 8, para product_type 16 (Quesos)
(17, 'A-01', 220, 20.0, 350, 35, 18.0, 9, 17), -- En warehouse 9, para product_type 17 (Conservas)
(18, 'B-01', 140, 4.0, 250, 25, 2.0, 9, 18), -- En warehouse 9, para product_type 18 (Snacks)
(19, 'A-01', 160, 15.0, 280, 28, 12.0, 10, 19), -- En warehouse 10, para product_type 19 (Pastas)
(20, 'B-01', 190, 10.0, 300, 30, 8.0, 10, 20); -- En warehouse 10, para product_type 20 (Salsas)

-- Insertando datos en 'products'
INSERT INTO `products` (`id`, `description`, `expiration_rate`, `freezing_rate`, `height`, `length`, `net_weight`, `product_code`, `recommended_freezing_temperature`, `width`, `product_type_id`, `seller_id`) VALUES
(1, 'Leche Entera 1L', 15, 0, 0.25, 0.1, 1.0, 'PROD-LE-01', 4.0, 0.1, 1, 1), -- Tipo: 1 (Lácteos), Vendedor: 1
(2, 'Carne de Res 1kg', 7, -18, 0.1, 0.2, 1.0, 'PROD-CR-01', -18.0, 0.15, 2, 2), -- Tipo: 2 (Carnes), Vendedor: 2
(3, 'Manzanas Royal Gala 1kg', 30, 0, 0.2, 0.3, 1.0, 'PROD-MA-01', 10.0, 0.2, 3, 3), -- Tipo: 3 (Frutas), Vendedor: 3
(4, 'Filete de Salmón 500g', 5, -20, 0.05, 0.25, 0.5, 'PROD-SA-01', -20.0, 0.15, 4, 4), -- Tipo: 4 (Pescados), Vendedor: 4
(5, 'Pizza Congelada Pepperoni', 365, -18, 0.04, 0.3, 0.5, 'PROD-PZ-01', -18.0, 0.3, 5, 5), -- Tipo: 5 (Congelados), Vendedor: 5
(6, 'Pechuga de Pollo 1kg', 6, -18, 0.1, 0.2, 1.0, 'PROD-PO-01', -18.0, 0.15, 6, 1), -- Tipo: 6 (Aves), Vendedor: 1
(7, 'Baguette Rústica', 2, 0, 0.08, 0.5, 0.4, 'PROD-PA-01', 20.0, 0.1, 7, 2), -- Tipo: 7 (Panadería), Vendedor: 2
(8, 'Refresco de Cola 2L', 180, 0, 0.3, 0.1, 2.0, 'PROD-RC-01', 15.0, 0.1, 8, 3), -- Tipo: 8 (Bebidas), Vendedor: 3
(9, 'Vino Tinto Malbec', 1825, 0, 0.3, 0.08, 0.75, 'PROD-VT-01', 14.0, 0.08, 9, 4), -- Tipo: 9 (Vinos), Vendedor: 4
(10, 'Helado de Vainilla 1L', 365, -22, 0.15, 0.15, 0.5, 'PROD-HV-01', -22.0, 0.1, 10, 5), -- Tipo: 10 (Helados), Vendedor: 5
(11, 'Lasaña Boloñesa Preparada', 10, -18, 0.06, 0.2, 0.4, 'PROD-LA-01', -18.0, 0.15, 11, 1), -- Tipo: 11 (Comidas P.), Vendedor: 1
(12, 'Tableta de Chocolate Negro', 365, 0, 0.01, 0.18, 0.1, 'PROD-CH-01', 20.0, 0.08, 12, 2), -- Tipo: 12 (Dulces), Vendedor: 2
(13, 'Café en Grano de Colombia 500g', 365, 0, 0.2, 0.1, 0.5, 'PROD-CA-01', 22.0, 0.08, 13, 3), -- Tipo: 13 (Café), Vendedor: 3
(14, 'Tomates Orgánicos 1kg', 10, 0, 0.15, 0.25, 1.0, 'PROD-TO-01', 12.0, 0.2, 14, 4), -- Tipo: 14 (Orgánicos), Vendedor: 4
(15, 'Queso Provolone 250g', 90, 0, 0.05, 0.1, 0.25, 'PROD-QP-01', 8.0, 0.1, 15, 5), -- Tipo: 15 (Verduras), Vendedor: 5
(16, 'Lata de Atún en Aceite', 1460, 0, 0.04, 0.08, 0.15, 'PROD-AT-01', 25.0, 0.08, 16, 1), -- Tipo: 16 (Quesos), Vendedor: 1
(17, 'Patatas Fritas sabor Queso', 180, 0, 0.3, 0.2, 0.2, 'PROD-PF-01', 20.0, 0.05, 17, 2), -- Tipo: 17 (Conservas), Vendedor: 2
(18, 'Pasta Fresca Tagliatelle 500g', 15, 0, 0.08, 0.2, 0.5, 'PROD-PT-01', 4.0, 0.15, 18, 3), -- Tipo: 18 (Snacks), Vendedor: 3
(19, 'Salsa de Tomate Casera 500g', 60, 0, 0.12, 0.08, 0.5, 'PROD-ST-01', 15.0, 0.08, 19, 4), -- Tipo: 19 (Pastas), Vendedor: 4
(20, 'Zanahorias 1kg', 20, 0, 0.1, 0.3, 1.0, 'PROD-ZA-01', 10.0, 0.15, 20, 5); -- Tipo: 20 (Salsas), Vendedor: 5

-- Insertando datos en 'product_records'
INSERT INTO `product_records` (`id`, `last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
(1, NOW(), 1.00, 1.50, 1), -- Registro de precio para product 1
(2, NOW(), 8.00, 12.00, 2), -- Registro de precio para product 2
(3, NOW(), 1.20, 2.00, 3), -- Registro de precio para product 3
(4, NOW(), 10.00, 15.00, 4), -- Registro de precio para product 4
(5, NOW(), 3.50, 5.00, 5), -- Registro de precio para product 5
(6, NOW(), 5.00, 7.50, 6), -- Registro de precio para product 6
(7, NOW(), 0.80, 1.20, 7), -- Registro de precio para product 7
(8, NOW(), 1.50, 2.20, 8), -- Registro de precio para product 8
(9, NOW(), 12.00, 20.00, 9), -- Registro de precio para product 9
(10, NOW(), 2.50, 4.00, 10), -- Registro de precio para product 10
(11, NOW(), 4.00, 6.00, 11), -- Registro de precio para product 11
(12, NOW(), 1.80, 2.50, 12), -- Registro de precio para product 12
(13, NOW(), 6.00, 10.00, 13), -- Registro de precio para product 13
(14, NOW(), 2.00, 3.50, 14), -- Registro de precio para product 14
(15, NOW(), 5.00, 8.00, 15), -- Registro de precio para product 15
(16, NOW(), 0.90, 1.40, 16), -- Registro de precio para product 16
(17, NOW(), 1.00, 1.60, 17), -- Registro de precio para product 17
(18, NOW(), 2.20, 3.80, 18), -- Registro de precio para product 18
(19, NOW(), 1.50, 2.50, 19), -- Registro de precio para product 19
(20, NOW(), 1.10, 1.90, 20); -- Registro de precio para product 20

-- Insertando datos en 'product_batches'
INSERT INTO `product_batches` (`id`, `batch_number`, `current_quantity`, `current_temperature`, `due_date`, `initial_quantity`, `manufacturing_date`, `manufacturing_hour`, `minimum_temperature`, `product_id`, `section_id`) VALUES
(1, 'L20250711-001', 1000, 4.0, '2025-07-26 00:00:00', 1000, '2025-07-11 00:00:00', '2025-07-11 08:00:00', 2.0, 1, 1), -- Lote de product 1 en section 1
(2, 'L20250710-002', 500, -18.0, '2025-07-17 00:00:00', 500, '2025-07-10 00:00:00', '2025-07-10 10:00:00', -22.0, 2, 2), -- Lote de product 2 en section 2
(3, 'L20250701-003', 2000, 10.0, '2025-07-31 00:00:00', 2000, '2025-07-01 00:00:00', '2025-07-01 06:00:00', 8.0, 3, 3), -- Lote de product 3 en section 3
(4, 'L20250712-004', 300, -20.0, '2025-07-17 00:00:00', 300, '2025-07-12 00:00:00', '2025-07-12 11:00:00', -25.0, 4, 4), -- Lote de product 4 en section 4
(5, 'L20250115-005', 1500, -18.0, '2026-01-15 00:00:00', 1500, '2025-01-15 00:00:00', '2025-01-15 14:00:00', -20.0, 5, 5), -- Lote de product 5 en section 5
(6, 'L20250711-006', 800, 2.0, '2025-07-17 00:00:00', 800, '2025-07-11 00:00:00', '2025-07-11 09:00:00', 0.0, 6, 6), -- Lote de product 6 en section 6
(7, 'L20250713-007', 500, 18.0, '2025-07-15 00:00:00', 500, '2025-07-13 00:00:00', '2025-07-13 05:00:00', 15.0, 7, 7), -- Lote de product 7 en section 7
(8, 'L20250620-008', 3000, 15.0, '2025-12-17 00:00:00', 3000, '2025-06-20 00:00:00', '2025-06-20 13:00:00', 12.0, 8, 8), -- Lote de product 8 en section 8
(9, 'L20230510-009', 600, 14.0, '2028-05-10 00:00:00', 600, '2023-05-10 00:00:00', '2023-05-10 18:00:00', 10.0, 9, 9), -- Lote de product 9 en section 9
(10, 'L20250201-010', 400, -22.0, '2026-02-01 00:00:00', 400, '2025-02-01 00:00:00', '2025-02-01 16:00:00', -25.0, 10, 10), -- Lote de product 10 en section 10
(11, 'L20250708-011', 250, -18.0, '2025-07-18 00:00:00', 250, '2025-07-08 00:00:00', '2025-07-08 12:00:00', -20.0, 11, 11), -- Lote de product 11 en section 11
(12, 'L20250415-012', 5000, 20.0, '2026-04-15 00:00:00', 5000, '2025-04-15 00:00:00', '2025-04-15 11:00:00', 18.0, 12, 12), -- Lote de product 12 en section 12
(13, 'L20250101-013', 1000, 22.0, '2026-01-01 00:00:00', 1000, '2025-01-01 00:00:00', '2025-01-01 07:00:00', 20.0, 13, 13), -- Lote de product 13 en section 13
(14, 'L20250709-014', 800, 12.0, '2025-07-19 00:00:00', 800, '2025-07-09 00:00:00', '2025-07-09 06:00:00', 10.0, 14, 14), -- Lote de product 14 en section 14
(15, 'L20250601-015', 300, 8.0, '2025-08-30 00:00:00', 300, '2025-06-01 00:00:00', '2025-06-01 10:00:00', 5.0, 15, 15), -- Lote de product 15 en section 15
(16, 'L20240820-016', 2000, 20.0, '2028-08-20 00:00:00', 2000, '2024-08-20 00:00:00', '2024-08-20 15:00:00', 18.0, 16, 16), -- Lote de product 16 en section 16
(17, 'L20250525-017', 4000, 20.0, '2025-11-21 00:00:00', 4000, '2025-05-25 00:00:00', '2025-05-25 17:00:00', 18.0, 17, 17), -- Lote de product 17 en section 17
(18, 'L20250710-018', 600, 4.0, '2025-07-25 00:00:00', 600, '2025-07-10 00:00:00', '2025-07-10 09:00:00', 2.0, 18, 18), -- Lote de product 18 en section 18
(19, 'L20250705-019', 350, 15.0, '2025-08-04 00:00:00', 350, '2025-07-05 00:00:00', '2025-07-05 14:00:00', 12.0, 19, 19), -- Lote de product 19 en section 19
(20, 'L20250703-020', 1200, 10.0, '2025-09-01 00:00:00', 1200, '2025-07-03 00:00:00', '2025-07-03 08:00:00', 8.0, 20, 20); -- Lote de product 20 en section 20

-- Insertando datos en 'inbound_orders'
INSERT INTO `inbound_orders` (`id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id`) VALUES
(1, NOW(), 'IN-2025-00001', 1, 1, 1), -- Recepción lote 1 por empleado 1 en warehouse 1
(2, NOW(), 'IN-2025-00002', 2, 2, 1), -- Recepción lote 2 por empleado 2 en warehouse 1
(3, NOW(), 'IN-2025-00003', 3, 3, 2), -- Recepción lote 3 por empleado 3 en warehouse 2
(4, NOW(), 'IN-2025-00004', 4, 4, 2), -- Recepción lote 4 por empleado 4 en warehouse 2
(5, NOW(), 'IN-2025-00005', 5, 5, 3), -- Recepción lote 5 por empleado 5 en warehouse 3
(6, NOW(), 'IN-2025-00006', 6, 6, 3), -- Recepción lote 6 por empleado 6 en warehouse 3
(7, NOW(), 'IN-2025-00007', 7, 7, 4), -- Recepción lote 7 por empleado 7 en warehouse 4
(8, NOW(), 'IN-2025-00008', 8, 8, 4), -- Recepción lote 8 por empleado 8 en warehouse 4
(9, NOW(), 'IN-2025-00009', 9, 9, 5), -- Recepción lote 9 por empleado 9 en warehouse 5
(10, NOW(), 'IN-2025-00010', 10, 10, 5), -- Recepción lote 10 por empleado 10 en warehouse 5
(11, NOW(), 'IN-2025-00011', 11, 11, 6), -- Recepción lote 11 por empleado 11 en warehouse 6
(12, NOW(), 'IN-2025-00012', 12, 12, 6), -- Recepción lote 12 por empleado 12 en warehouse 6
(13, NOW(), 'IN-2025-00013', 13, 13, 7), -- Recepción lote 13 por empleado 13 en warehouse 7
(14, NOW(), 'IN-2025-00014', 14, 14, 7), -- Recepción lote 14 por empleado 14 en warehouse 7
(15, NOW(), 'IN-2025-00015', 15, 15, 8), -- Recepción lote 15 por empleado 15 en warehouse 8
(16, NOW(), 'IN-2025-00016', 16, 16, 8), -- Recepción lote 16 por empleado 16 en warehouse 8
(17, NOW(), 'IN-2025-00017', 17, 17, 9), -- Recepción lote 17 por empleado 17 en warehouse 9
(18, NOW(), 'IN-2025-00018', 18, 18, 9), -- Recepción lote 18 por empleado 18 en warehouse 9
(19, NOW(), 'IN-2025-00019', 19, 19, 10), -- Recepción lote 19 por empleado 19 en warehouse 10
(20, NOW(), 'IN-2025-00020', 20, 20, 10); -- Recepción lote 20 por empleado 20 en warehouse 10

-- Insertando datos en 'carriers'
INSERT INTO `carriers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES
(1, 'CAR-001', 'Transportes Veloz', 'Calle Rápida 456', '555-0301', 1), -- Ubicado en locality 1 (Bogotá)
(2, 'CAR-002', 'Logística Global', 'Avenida Mundo 12', '555-0302', 2), -- Ubicado en locality 2 (Medellín)
(3, 'CAR-003', 'Envíos Seguros S.A.', 'Carrera Confianza 88', '555-0303', 3), -- Ubicado en locality 3 (Cali)
(4, 'CAR-004', 'Carga Fría Express', 'Ruta Helada 99', '555-0304', 4), -- Ubicado en locality 4 (La Plata)
(5, 'CAR-005', 'Distribución Andina', 'Cordillera 101', '555-0305', 5), -- Ubicado en locality 5 (Córdoba)
(6, 'CAR-006', 'Trans-Pacífico', 'Costa 202', '555-0306', 6), -- Ubicado en locality 6 (Rosario)
(7, 'CAR-007', 'Logística del Golfo', 'Bahía 303', '555-0307', 7), -- Ubicado en locality 7 (Guadalajara)
(8, 'CAR-008', 'EuroTrans', 'Continente 404', '555-0308', 8), -- Ubicado en locality 8 (Monterrey)
(9, 'CAR-009', 'Iberia Cargo', 'Península 505', '555-0309', 9), -- Ubicado en locality 9 (CDMX)
(10, 'CAR-010', 'MercoSur Logística', 'Mercado 606', '555-0310', 1), -- Ubicado en locality 1 (Bogotá)
(11, 'CAR-011', 'Transportes del Sur', 'Cono Sur 707', '555-0311', 2), -- Ubicado en locality 2 (Medellín)
(12, 'CAR-012', 'Amazonia Cargo', 'Selva 808', '555-0312', 3), -- Ubicado en locality 3 (Cali)
(13, 'CAR-013', 'Caribe Envíos', 'Mar Caribe 909', '555-0313', 4), -- Ubicado en locality 4 (La Plata)
(14, 'CAR-014', 'Norteamérica Freight', 'Ruta 66', '555-0314', 5), -- Ubicado en locality 5 (Córdoba)
(15, 'CAR-015', 'Canada Connect', 'Trans-Canada Hwy 1', '555-0315', 6), -- Ubicado en locality 6 (Rosario)
(16, 'CAR-016', 'La Poste Cargo', 'Rue de la Logistique 20', '555-0316', 7), -- Ubicado en locality 7 (Guadalajara)
(17, 'CAR-017', 'DHL', 'Global Avenue 1', '555-0317', 8), -- Ubicado en locality 8 (Monterrey)
(18, 'CAR-018', 'FedEx', 'Express Lane 2', '555-0318', 9), -- Ubicado en locality 9 (CDMX)
(19, 'CAR-019', 'UPS', 'Worldwide Service 3', '555-0319', 1), -- Ubicado en locality 1 (Bogotá)
(20, 'CAR-020', 'Servientrega', 'Diagonal 100', '555-0320', 2); -- Ubicado en locality 2 (Medellín)

-- Insertando datos en 'order_status'
INSERT INTO `order_status` (`id`, `description`) VALUES
(1, 'Procesando'), (2, 'Confirmado'), (3, 'Preparando Envío'), (4, 'Enviado'),
(5, 'En Tránsito'), (6, 'Entregado'), (7, 'Cancelado'), (8, 'Devuelto'),
(9, 'En espera de pago'), (10, 'Pago recibido'), (11, 'En espera de stock'), (12, 'Pedido parcial'),
(13, 'Error en pedido'), (14, 'Revisión manual'), (15, 'Listo para recoger'), (16, 'Recogido por transportista'),
(17, 'En aduanas'), (18, 'Retrasado'), (19, 'Completado'), (20, 'Cerrado');

-- Insertando datos en 'purchase_orders'
INSERT INTO `purchase_orders` (`id`, `order_number`, `order_date`, `tracking_code`, `buyer_id`, `product_record_id`) VALUES
(1, 'PO-2025-00001', NOW(), 'TRK123456789', 1, 1), -- Comprador 1 (Ignacio) compra product_record 1
(2, 'PO-2025-00002', NOW(), 'TRK123456790', 2, 2), -- Comprador 2 (Julian) compra product_record 2
(3, 'PO-2025-00003', NOW(), 'TRK123456791', 3, 3), -- Comprador 3 (Karen) compra product_record 3
(4, 'PO-2025-00004', NOW(), 'TRK123456792', 4, 4), -- Comprador 4 (Samuel) compra product_record 4
(5, 'PO-2025-00005', NOW(), 'TRK123456793', 5, 5), -- Comprador 5 (Gabriel) compra product_record 5
(6, 'PO-2025-00006', NOW(), 'TRK123456794', 6, 6), -- Comprador 6 (Juan) compra product_record 6
(7, 'PO-2025-00007', NOW(), 'TRK123456795', 1, 7), -- Comprador 1 (Ignacio) compra product_record 7
(8, 'PO-2025-00008', NOW(), 'TRK123456796', 2, 8), -- Comprador 2 (Julian) compra product_record 8
(9, 'PO-2025-00009', NOW(), 'TRK123456797', 3, 9), -- Comprador 3 (Karen) compra product_record 9
(10, 'PO-2025-00010', NOW(), 'TRK123456798', 4, 10), -- Comprador 4 (Samuel) compra product_record 10
(11, 'PO-2025-00011', NOW(), 'TRK123456799', 5, 11), -- Comprador 5 (Gabriel) compra product_record 11
(12, 'PO-2025-00012', NOW(), 'TRK123456800', 6, 12), -- Comprador 6 (Juan) compra product_record 12
(13, 'PO-2025-00013', NOW(), 'TRK123456801', 1, 13), -- Comprador 1 (Ignacio) compra product_record 13
(14, 'PO-2025-00014', NOW(), 'TRK123456802', 2, 14), -- Comprador 2 (Julian) compra product_record 14
(15, 'PO-2025-00015', NOW(), 'TRK123456803', 3, 15), -- Comprador 3 (Karen) compra product_record 15
(16, 'PO-2025-00016', NOW(), 'TRK123456804', 4, 16), -- Comprador 4 (Samuel) compra product_record 16
(17, 'PO-2025-00017', NOW(), 'TRK123456805', 5, 17), -- Comprador 5 (Gabriel) compra product_record 17
(18, 'PO-2025-00018', NOW(), 'TRK123456806', 6, 18), -- Comprador 6 (Juan) compra product_record 18
(19, 'PO-2025-00019', NOW(), 'TRK123456807', 1, 19), -- Comprador 1 (Ignacio) compra product_record 19
(20, 'PO-2025-00020', NOW(), 'TRK123456808', 2, 20); -- Comprador 2 (Julian) compra product_record 20

-- Insertando datos en 'users'
INSERT INTO `users` (`id`, `username`, `password`) VALUES
(1, 'admin', 'hashed_password_1'), (2, 'jperez', 'hashed_password_2'), (3, 'mgarcia', 'hashed_password_3'),
(4, 'pramirez', 'hashed_password_4'), (5, 'lfernandez', 'hashed_password_5'), (6, 'alopez', 'hashed_password_6'),
(7, 'csanz', 'hashed_password_7'), (8, 'dmoreno', 'hashed_password_8'), (9, 'bjimenez', 'hashed_password_9'),
(10, 'sruiz', 'hashed_password_10'), (11, 'ralonso', 'hashed_password_11'), (12, 'mgutierrez', 'hashed_password_12'),
(13, 'enavarro', 'hashed_password_13'), (14, 'jiglesias', 'hashed_password_14'), (15, 'cblanco', 'hashed_password_15'),
(16, 'rsoto', 'hashed_password_16'), (17, 'ncrespo', 'hashed_password_17'), (18, 'freyes', 'hashed_password_18'),
(19, 'vgil', 'hashed_password_19'), (20, 'oscaro', 'hashed_password_20');

-- Insertando datos en 'rol'
INSERT INTO `rol` (`id`, `rol_name`, `description`) VALUES
(1, 'Administrador', 'Acceso total al sistema'), (2, 'Gerente de Almacén', 'Gestiona un almacén específico'),
(3, 'Operario de Almacén', 'Realiza tareas de recepción y despacho'), (4, 'Cliente', 'Acceso para ver y realizar pedidos'),
(5, 'Analista de Datos', 'Acceso para generar reportes'), (6, 'Soporte Técnico', 'Mantenimiento del sistema'),
(7, 'Auditor', 'Revisa la integridad de los datos'), (8, 'Gerente de Compras', 'Autoriza órdenes de compra'),
(9, 'Vendedor', 'Gestiona clientes y ventas'), (10, 'Contador', 'Acceso a datos financieros'),
(11, 'Recursos Humanos', 'Gestiona empleados'), (12, 'Supervisor de Calidad', 'Control de calidad de productos'),
(13, 'Jefe de Logística', 'Coordina transportistas y rutas'), (14, 'Asistente Administrativo', 'Apoyo en tareas administrativas'),
(15, 'Marketing', 'Análisis de mercado y productos'), (16, 'Proveedor', 'Acceso para ver estado de sus productos'),
(17, 'Director General', 'Vista global de la operación'), (18, 'Pasante', 'Acceso limitado para aprendizaje'),
(19, 'Seguridad', 'Monitoreo de accesos'), (20, 'Invitado', 'Acceso de solo lectura muy limitado');

-- Insertando datos en 'user_rol'
INSERT INTO `user_rol` (`id`, `usuario_id`, `rol_id`) VALUES
(1, 1, 1), -- usuario 1 (admin) tiene rol 1 (Administrador)
(2, 2, 3), -- usuario 2 (jperez) tiene rol 3 (Operario)
(3, 3, 3), -- usuario 3 (mgarcia) tiene rol 3 (Operario)
(4, 4, 3), -- usuario 4 (pramirez) tiene rol 3 (Operario)
(5, 5, 2), -- usuario 5 (lfernandez) tiene rol 2 (Gerente)
(6, 6, 3), -- usuario 6 (alopez) tiene rol 3 (Operario)
(7, 7, 3), -- usuario 7 (csanz) tiene rol 3 (Operario)
(8, 8, 2), -- usuario 8 (dmoreno) tiene rol 2 (Gerente)
(9, 9, 3), -- usuario 9 (bjimenez) tiene rol 3 (Operario)
(10, 10, 2), -- usuario 10 (sruiz) tiene rol 2 (Gerente)
(11, 11, 3), -- usuario 11 (ralonso) tiene rol 3 (Operario)
(12, 12, 3), -- usuario 12 (mgutierrez) tiene rol 3 (Operario)
(13, 13, 2), -- usuario 13 (enavarro) tiene rol 2 (Gerente)
(14, 14, 2), -- usuario 14 (jiglesias) tiene rol 2 (Gerente)
(15, 15, 2), -- usuario 15 (cblanco) tiene rol 2 (Gerente)
(16, 16, 13), -- usuario 16 (rsoto) tiene rol 13 (Logística)
(17, 17, 13), -- usuario 17 (ncrespo) tiene rol 13 (Logística)
(18, 18, 13), -- usuario 18 (freyes) tiene rol 13 (Logística)
(19, 19, 12), -- usuario 19 (vgil) tiene rol 12 (Calidad)
(20, 20, 12); -- usuario 20 (oscaro) tiene rol 12 (Calidad)