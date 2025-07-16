-- =================================================================
-- SCRIPT DE CREACIÓN DE LA BASE DE DATOS Y TABLAS
-- =================================================================

DROP DATABASE IF EXISTS productos_frescos;

CREATE DATABASE productos_frescos;

USE productos_frescos;

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
-- Si se elimina un país, se eliminarán todas sus provincias.
CREATE TABLE `provinces` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `province_name` VARCHAR(255) NOT NULL,
  `id_country_fk` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`id_country_fk`) REFERENCES `countries`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'localities'
-- Si se elimina una provincia, se eliminarán todas sus localidades.
CREATE TABLE `localities` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `locality_name` VARCHAR(255) NOT NULL,
  `province_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`province_id`) REFERENCES `provinces`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'sellers'
-- Si se elimina una localidad, se eliminarán los vendedores asociados.
CREATE TABLE `sellers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `cid` VARCHAR(255) NOT NULL,
  `company_name` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `telephone` VARCHAR(255) NOT NULL,
  `locality_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`locality_id`) REFERENCES `localities`(`id`) ON DELETE CASCADE
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
-- Si se elimina una localidad, se eliminarán los almacenes asociados.
CREATE TABLE `warehouse` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `address` VARCHAR(255) NOT NULL,
  `telephone` VARCHAR(255) NOT NULL,
  `minimum_temperature` INT,
  `minimum_capacity` INT,
  `warehouse_code` VARCHAR(255) NOT NULL,
  `locality_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`locality_id`) REFERENCES `localities`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'employees'
-- Si se elimina un almacén, se eliminarán los empleados asociados.
CREATE TABLE `employees` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `id_card_number` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(255) NOT NULL,
  `last_name` VARCHAR(255) NOT NULL,
  `warehouse_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'products_types'
CREATE TABLE `products_types` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'sections'
-- Si se elimina un almacén o un tipo de producto, se eliminarán las secciones relacionadas.
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
  FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`product_type_id`) REFERENCES `products_types`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'products'
-- Si se elimina un tipo de producto o un vendedor, se eliminarán los productos asociados.
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
  FOREIGN KEY (`product_type_id`) REFERENCES `products_types`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`seller_id`) REFERENCES `sellers`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'product_records'
-- Si se elimina un producto, se eliminarán sus registros.
CREATE TABLE `product_records` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `last_update_date` DATETIME(6),
  `purchase_price` DECIMAL(19,2),
  `sale_price` DECIMAL(19,2),
  `product_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'product_batches'
-- Si se elimina un producto o una sección, se eliminarán los lotes de productos asociados.
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
  FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`section_id`) REFERENCES `sections`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'inbound_orders'
-- Si se elimina un empleado, lote de producto o almacén, se eliminarán las órdenes de entrada relacionadas.
CREATE TABLE `inbound_orders` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_date` DATETIME(6),
  `order_number` VARCHAR(255) NOT NULL,
  `employee_id` INT NOT NULL,
  `product_batch_id` INT NOT NULL,
  `warehouse_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`employee_id`) REFERENCES `employees`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`product_batch_id`) REFERENCES `product_batches`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'carriers'
-- Si se elimina una localidad, se eliminarán los transportistas asociados.
CREATE TABLE `carriers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `cid` VARCHAR(255) NOT NULL,
  `company_name` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `telephone` VARCHAR(255) NOT NULL,
  `locality_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`locality_id`) REFERENCES `localities`(`id`) ON DELETE CASCADE
);

-- Creación de la tabla 'order_status'
CREATE TABLE `order_status` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(255),
  PRIMARY KEY (`id`)
);

-- Creación de la tabla 'purchase_orders'
-- Si se elimina un comprador o un registro de producto, se eliminarán las órdenes de compra asociadas.
CREATE TABLE `purchase_orders` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_number` VARCHAR(255) NOT NULL,
  `order_date` DATETIME NOT NULL,
  `tracking_code` VARCHAR(255),
  `buyer_id` INT NOT NULL,
  `product_record_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`buyer_id`) REFERENCES `buyers`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`product_record_id`) REFERENCES `product_records`(`id`) ON DELETE CASCADE
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
-- Si se elimina un usuario o un rol, se eliminará la relación en esta tabla.
CREATE TABLE `user_rol` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `usuario_id` INT NOT NULL,
  `rol_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`usuario_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`rol_id`) REFERENCES `rol`(`id`) ON DELETE CASCADE
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

-- =================================================================
-- SCRIPT DE INSERCIÓN DE DATOS PARA 'product_records'
-- Se generan entre 10 y 20 registros por cada producto para simular un historial de precios.
-- =================================================================

-- --- Producto 1: Leche Entera 1L (15 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-05-10 10:00:00', 1.00, 1.50, 1), ('2024-06-10 10:00:00', 1.02, 1.55, 1),
('2024-07-10 10:00:00', 1.03, 1.55, 1), ('2024-08-10 10:00:00', 1.02, 1.53, 1),
('2024-09-10 10:00:00', 1.04, 1.56, 1), ('2024-10-10 10:00:00', 1.05, 1.58, 1),
('2024-11-10 10:00:00', 1.05, 1.60, 1), ('2024-12-10 10:00:00', 1.06, 1.60, 1),
('2025-01-10 10:00:00', 1.07, 1.62, 1), ('2025-02-10 10:00:00', 1.08, 1.65, 1),
('2025-03-10 10:00:00', 1.07, 1.64, 1), ('2025-04-10 10:00:00', 1.08, 1.65, 1),
('2025-05-10 10:00:00', 1.09, 1.68, 1), ('2025-06-10 10:00:00', 1.10, 1.70, 1),
(NOW(), 1.12, 1.72, 1);

-- --- Producto 2: Carne de Res 1kg (12 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-08-15 11:00:00', 8.15, 12.25, 2), ('2024-09-15 11:00:00', 8.25, 12.40, 2),
('2024-10-15 11:00:00', 8.30, 12.50, 2), ('2024-11-15 11:00:00', 8.35, 12.60, 2),
('2024-12-15 11:00:00', 8.40, 12.70, 2), ('2025-01-15 11:00:00', 8.45, 12.80, 2),
('2025-02-15 11:00:00', 8.50, 12.90, 2), ('2025-03-15 11:00:00', 8.45, 12.85, 2),
('2025-04-15 11:00:00', 8.55, 13.00, 2), ('2025-05-15 11:00:00', 8.60, 13.10, 2),
('2025-06-15 11:00:00', 8.65, 13.20, 2), (NOW(), 8.75, 13.40, 2);

-- --- Producto 3: Manzanas Royal Gala 1kg (18 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-02-20 09:30:00', 1.19, 1.98, 3), ('2024-03-20 09:30:00', 1.20, 2.00, 3),
('2024-04-20 09:30:00', 1.21, 2.02, 3), ('2024-05-20 09:30:00', 1.20, 2.00, 3),
('2024-06-20 09:30:00', 1.22, 2.05, 3), ('2024-07-20 09:30:00', 1.23, 2.05, 3),
('2024-08-20 09:30:00', 1.22, 2.03, 3), ('2024-09-20 09:30:00', 1.24, 2.06, 3),
('2024-10-20 09:30:00', 1.25, 2.08, 3), ('2024-11-20 09:30:00', 1.25, 2.10, 3),
('2024-12-20 09:30:00', 1.26, 2.10, 3), ('2025-01-20 09:30:00', 1.27, 2.12, 3),
('2025-02-20 09:30:00', 1.28, 2.15, 3), ('2025-03-20 09:30:00', 1.27, 2.14, 3),
('2025-04-20 09:30:00', 1.28, 2.15, 3), ('2025-05-20 09:30:00', 1.29, 2.18, 3),
('2025-06-20 09:30:00', 1.30, 2.20, 3), (NOW(), 1.32, 2.22, 3);

-- --- Producto 4: Filete de Salmón 500g (11 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-09-25 14:00:00', 10.25, 15.40, 4), ('2024-10-25 14:00:00', 10.30, 15.50, 4),
('2024-11-25 14:00:00', 10.35, 15.60, 4), ('2024-12-25 14:00:00', 10.40, 15.70, 4),
('2025-01-25 14:00:00', 10.45, 15.80, 4), ('2025-02-25 14:00:00', 10.50, 15.90, 4),
('2025-03-25 14:00:00', 10.45, 15.85, 4), ('2025-04-25 14:00:00', 10.55, 16.00, 4),
('2025-05-25 14:00:00', 10.60, 16.10, 4), ('2025-06-25 14:00:00', 10.65, 16.20, 4),
(NOW(), 10.75, 16.40, 4);

-- --- Producto 5: Pizza Congelada Pepperoni (20 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-01-01 12:00:00', 3.45, 4.90, 5), ('2024-02-01 12:00:00', 3.48, 4.95, 5),
('2024-03-01 12:00:00', 3.50, 5.00, 5), ('2024-04-01 12:00:00', 3.52, 5.05, 5),
('2024-05-01 12:00:00', 3.50, 5.00, 5), ('2024-06-01 12:00:00', 3.55, 5.10, 5),
('2024-07-01 12:00:00', 3.58, 5.15, 5), ('2024-08-01 12:00:00', 3.56, 5.12, 5),
('2024-09-01 12:00:00', 3.60, 5.20, 5), ('2024-10-01 12:00:00', 3.62, 5.25, 5),
('2024-11-01 12:00:00', 3.65, 5.30, 5), ('2024-12-01 12:00:00', 3.68, 5.35, 5),
('2025-01-01 12:00:00', 3.70, 5.40, 5), ('2025-02-01 12:00:00', 3.72, 5.45, 5),
('2025-03-01 12:00:00', 3.70, 5.42, 5), ('2025-04-01 12:00:00', 3.75, 5.50, 5),
('2025-05-01 12:00:00', 3.78, 5.55, 5), ('2025-06-01 12:00:00', 3.80, 5.60, 5),
('2025-07-01 12:00:00', 3.82, 5.65, 5), (NOW(), 3.85, 5.70, 5);

-- --- Producto 6: Pechuga de Pollo 1kg (14 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-06-05 15:00:00', 5.10, 7.60, 6), ('2024-07-05 15:00:00', 5.15, 7.65, 6),
('2024-08-05 15:00:00', 5.10, 7.60, 6), ('2024-09-05 15:00:00', 5.20, 7.70, 6),
('2024-10-05 15:00:00', 5.25, 7.75, 6), ('2024-11-05 15:00:00', 5.30, 7.80, 6),
('2024-12-05 15:00:00', 5.35, 7.85, 6), ('2025-01-05 15:00:00', 5.40, 7.90, 6),
('2025-02-05 15:00:00', 5.45, 7.95, 6), ('2025-03-05 15:00:00', 5.40, 7.90, 6),
('2025-04-05 15:00:00', 5.50, 8.00, 6), ('2025-05-05 15:00:00', 5.55, 8.05, 6),
('2025-06-05 15:00:00', 5.60, 8.10, 6), (NOW(), 5.70, 8.20, 6);

-- --- Producto 7: Baguette Rústica (17 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-03-08 07:00:00', 0.80, 1.20, 7), ('2024-04-08 07:00:00', 0.81, 1.21, 7),
('2024-05-08 07:00:00', 0.80, 1.20, 7), ('2024-06-08 07:00:00', 0.82, 1.22, 7),
('2024-07-08 07:00:00', 0.83, 1.23, 7), ('2024-08-08 07:00:00', 0.82, 1.22, 7),
('2024-09-08 07:00:00', 0.84, 1.24, 7), ('2024-10-08 07:00:00', 0.85, 1.25, 7),
('2024-11-08 07:00:00', 0.85, 1.25, 7), ('2024-12-08 07:00:00', 0.86, 1.26, 7),
('2025-01-08 07:00:00', 0.87, 1.27, 7), ('2025-02-08 07:00:00', 0.88, 1.28, 7),
('2025-03-08 07:00:00', 0.87, 1.27, 7), ('2025-04-08 07:00:00', 0.88, 1.28, 7),
('2025-05-08 07:00:00', 0.89, 1.29, 7), ('2025-06-08 07:00:00', 0.90, 1.30, 7),
(NOW(), 0.92, 1.32, 7);

-- --- Producto 8: Refresco de Cola 2L (10 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-10-12 16:00:00', 1.55, 2.25, 8), ('2024-11-12 16:00:00', 1.55, 2.25, 8),
('2024-12-12 16:00:00', 1.56, 2.26, 8), ('2025-01-12 16:00:00', 1.57, 2.27, 8),
('2025-02-12 16:00:00', 1.58, 2.28, 8), ('2025-03-12 16:00:00', 1.57, 2.27, 8),
('2025-04-12 16:00:00', 1.58, 2.28, 8), ('2025-05-12 16:00:00', 1.59, 2.29, 8),
('2025-06-12 16:00:00', 1.60, 2.30, 8), (NOW(), 1.62, 2.32, 8);

-- --- Producto 9: Vino Tinto Malbec (19 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-02-19 18:00:00', 11.90, 19.90, 9), ('2024-03-19 18:00:00', 12.00, 20.00, 9),
('2024-04-19 18:00:00', 12.10, 20.10, 9), ('2024-05-19 18:00:00', 12.05, 20.05, 9),
('2024-06-19 18:00:00', 12.15, 20.20, 9), ('2024-07-19 18:00:00', 12.20, 20.30, 9),
('2024-08-19 18:00:00', 12.15, 20.25, 9), ('2024-09-19 18:00:00', 12.25, 20.40, 9),
('2024-10-19 18:00:00', 12.30, 20.50, 9), ('2024-11-19 18:00:00', 12.35, 20.60, 9),
('2024-12-19 18:00:00', 12.40, 20.70, 9), ('2025-01-19 18:00:00', 12.45, 20.80, 9),
('2025-02-19 18:00:00', 12.50, 20.90, 9), ('2025-03-19 18:00:00', 12.45, 20.85, 9),
('2025-04-19 18:00:00', 12.55, 21.00, 9), ('2025-05-19 18:00:00', 12.60, 21.10, 9),
('2025-06-19 18:00:00', 12.65, 21.20, 9), ('2025-07-19 18:00:00', 12.70, 21.30, 9),
(NOW(), 12.75, 21.40, 9);

-- --- Producto 10: Helado de Vainilla 1L (13 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-07-22 13:00:00', 2.58, 4.15, 10), ('2024-08-22 13:00:00', 2.56, 4.12, 10),
('2024-09-22 13:00:00', 2.60, 4.20, 10), ('2024-10-22 13:00:00', 2.62, 4.25, 10),
('2024-11-22 13:00:00', 2.65, 4.30, 10), ('2024-12-22 13:00:00', 2.68, 4.35, 10),
('2025-01-22 13:00:00', 2.70, 4.40, 10), ('2025-02-22 13:00:00', 2.72, 4.45, 10),
('2025-03-22 13:00:00', 2.70, 4.42, 10), ('2025-04-22 13:00:00', 2.75, 4.50, 10),
('2025-05-22 13:00:00', 2.78, 4.55, 10), ('2025-06-22 13:00:00', 2.80, 4.60, 10),
(NOW(), 2.85, 4.70, 10);

-- --- Producto 11: Lasaña Boloñesa Preparada (16 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-04-28 17:00:00', 4.05, 6.05, 11), ('2024-05-28 17:00:00', 4.00, 6.00, 11),
('2024-06-28 17:00:00', 4.10, 6.10, 11), ('2024-07-28 17:00:00', 4.15, 6.15, 11),
('2024-08-28 17:00:00', 4.10, 6.10, 11), ('2024-09-28 17:00:00', 4.20, 6.20, 11),
('2024-10-28 17:00:00', 4.25, 6.25, 11), ('2024-11-28 17:00:00', 4.30, 6.30, 11),
('2024-12-28 17:00:00', 4.35, 6.35, 11), ('2025-01-28 17:00:00', 4.40, 6.40, 11),
('2025-02-28 17:00:00', 4.45, 6.45, 11), ('2025-03-28 17:00:00', 4.40, 6.40, 11),
('2025-04-28 17:00:00', 4.50, 6.50, 11), ('2025-05-28 17:00:00', 4.55, 6.55, 11),
('2025-06-28 17:00:00', 4.60, 6.60, 11), (NOW(), 4.70, 6.70, 11);

-- --- Producto 12: Tableta de Chocolate Negro (15 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-05-14 19:00:00', 1.80, 2.50, 12), ('2024-06-14 19:00:00', 1.85, 2.55, 12),
('2024-07-14 19:00:00', 1.88, 2.58, 12), ('2024-08-14 19:00:00', 1.86, 2.56, 12),
('2024-09-14 19:00:00', 1.90, 2.60, 12), ('2024-10-14 19:00:00', 1.92, 2.62, 12),
('2024-11-14 19:00:00', 1.95, 2.65, 12), ('2024-12-14 19:00:00', 1.98, 2.68, 12),
('2025-01-14 19:00:00', 2.00, 2.70, 12), ('2025-02-14 19:00:00', 2.02, 2.72, 12),
('2025-03-14 19:00:00', 2.00, 2.70, 12), ('2025-04-14 19:00:00', 2.05, 2.75, 12),
('2025-05-14 19:00:00', 2.08, 2.78, 12), ('2025-06-14 19:00:00', 2.10, 2.80, 12),
(NOW(), 2.15, 2.85, 12);

-- --- Producto 13: Café en Grano de Colombia 500g (12 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-08-16 06:00:00', 6.10, 10.10, 13), ('2024-09-16 06:00:00', 6.20, 10.20, 13),
('2024-10-16 06:00:00', 6.25, 10.25, 13), ('2024-11-16 06:00:00', 6.30, 10.30, 13),
('2024-12-16 06:00:00', 6.35, 10.35, 13), ('2025-01-16 06:00:00', 6.40, 10.40, 13),
('2025-02-16 06:00:00', 6.45, 10.45, 13), ('2025-03-16 06:00:00', 6.40, 10.40, 13),
('2025-04-16 06:00:00', 6.50, 10.50, 13), ('2025-05-16 06:00:00', 6.55, 10.55, 13),
('2025-06-16 06:00:00', 6.60, 10.60, 13), (NOW(), 6.70, 10.70, 13);

-- --- Producto 14: Tomates Orgánicos 1kg (18 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-02-21 08:30:00', 1.98, 3.48, 14), ('2024-03-21 08:30:00', 2.00, 3.50, 14),
('2024-04-21 08:30:00', 2.02, 3.52, 14), ('2024-05-21 08:30:00', 2.00, 3.50, 14),
('2024-06-21 08:30:00', 2.05, 3.55, 14), ('2024-07-21 08:30:00', 2.08, 3.58, 14),
('2024-08-21 08:30:00', 2.06, 3.56, 14), ('2024-09-21 08:30:00', 2.10, 3.60, 14),
('2024-10-21 08:30:00', 2.12, 3.62, 14), ('2024-11-21 08:30:00', 2.15, 3.65, 14),
('2024-12-21 08:30:00', 2.18, 3.68, 14), ('2025-01-21 08:30:00', 2.20, 3.70, 14),
('2025-02-21 08:30:00', 2.22, 3.72, 14), ('2025-03-21 08:30:00', 2.20, 3.70, 14),
('2025-04-21 08:30:00', 2.25, 3.75, 14), ('2025-05-21 08:30:00', 2.28, 3.78, 14),
('2025-06-21 08:30:00', 2.30, 3.80, 14), (NOW(), 2.35, 3.85, 14);

-- --- Producto 15: Queso Provolone 250g (11 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-09-26 11:30:00', 5.20, 8.20, 15), ('2024-10-26 11:30:00', 5.25, 8.25, 15),
('2024-11-26 11:30:00', 5.30, 8.30, 15), ('2024-12-26 11:30:00', 5.35, 8.35, 15),
('2025-01-26 11:30:00', 5.40, 8.40, 15), ('2025-02-26 11:30:00', 5.45, 8.45, 15),
('2025-03-26 11:30:00', 5.40, 8.40, 15), ('2025-04-26 11:30:00', 5.50, 8.50, 15),
('2025-05-26 11:30:00', 5.55, 8.55, 15), ('2025-06-26 11:30:00', 5.60, 8.60, 15),
(NOW(), 5.70, 8.70, 15);

-- --- Producto 16: Lata de Atún en Aceite (20 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-01-30 10:30:00', 0.88, 1.38, 16), ('2024-02-28 10:30:00', 0.89, 1.39, 16),
('2024-03-30 10:30:00', 0.90, 1.40, 16), ('2024-04-30 10:30:00', 0.91, 1.41, 16),
('2024-05-30 10:30:00', 0.90, 1.40, 16), ('2024-06-30 10:30:00', 0.92, 1.42, 16),
('2024-07-30 10:30:00', 0.93, 1.43, 16), ('2024-08-30 10:30:00', 0.92, 1.42, 16),
('2024-09-30 10:30:00', 0.94, 1.44, 16), ('2024-10-30 10:30:00', 0.95, 1.45, 16),
('2024-11-30 10:30:00', 0.95, 1.45, 16), ('2024-12-30 10:30:00', 0.96, 1.46, 16),
('2025-01-30 10:30:00', 0.97, 1.47, 16), ('2025-02-28 10:30:00', 0.98, 1.48, 16),
('2025-03-30 10:30:00', 0.97, 1.47, 16), ('2025-04-30 10:30:00', 0.98, 1.48, 16),
('2025-05-30 10:30:00', 0.99, 1.49, 16), ('2025-06-30 10:30:00', 1.00, 1.50, 16),
('2025-07-30 10:30:00', 1.00, 1.50, 16), (NOW(), 1.02, 1.52, 16);

-- --- Producto 17: Patatas Fritas sabor Queso (14 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-06-04 20:00:00', 1.02, 1.62, 17), ('2024-07-04 20:00:00', 1.03, 1.63, 17),
('2024-08-04 20:00:00', 1.02, 1.62, 17), ('2024-09-04 20:00:00', 1.04, 1.64, 17),
('2024-10-04 20:00:00', 1.05, 1.65, 17), ('2024-11-04 20:00:00', 1.05, 1.65, 17),
('2024-12-04 20:00:00', 1.06, 1.66, 17), ('2025-01-04 20:00:00', 1.07, 1.67, 17),
('2025-02-04 20:00:00', 1.08, 1.68, 17), ('2025-03-04 20:00:00', 1.07, 1.67, 17),
('2025-04-04 20:00:00', 1.08, 1.68, 17), ('2025-05-04 20:00:00', 1.09, 1.69, 17),
('2025-06-04 20:00:00', 1.10, 1.70, 17), (NOW(), 1.12, 1.72, 17);

-- --- Producto 18: Pasta Fresca Tagliatelle 500g (17 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-03-07 12:30:00', 2.20, 3.80, 18), ('2024-04-07 12:30:00', 2.22, 3.82, 18),
('2024-05-07 12:30:00', 2.20, 3.80, 18), ('2024-06-07 12:30:00', 2.25, 3.85, 18),
('2024-07-07 12:30:00', 2.28, 3.88, 18), ('2024-08-07 12:30:00', 2.26, 3.86, 18),
('2024-09-07 12:30:00', 2.30, 3.90, 18), ('2024-10-07 12:30:00', 2.32, 3.92, 18),
('2024-11-07 12:30:00', 2.35, 3.95, 18), ('2024-12-07 12:30:00', 2.38, 3.98, 18),
('2025-01-07 12:30:00', 2.40, 4.00, 18), ('2025-02-07 12:30:00', 2.42, 4.02, 18),
('2025-03-07 12:30:00', 2.40, 4.00, 18), ('2025-04-07 12:30:00', 2.45, 4.05, 18),
('2025-05-07 12:30:00', 2.48, 4.08, 18), ('2025-06-07 12:30:00', 2.50, 4.10, 18),
(NOW(), 2.55, 4.15, 18);

-- --- Producto 19: Salsa de Tomate Casera 500g (10 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-10-11 10:00:00', 1.62, 2.62, 19), ('2024-11-11 10:00:00', 1.65, 2.65, 19),
('2024-12-11 10:00:00', 1.68, 2.68, 19), ('2025-01-11 10:00:00', 1.70, 2.70, 19),
('2025-02-11 10:00:00', 1.72, 2.72, 19), ('2025-03-11 10:00:00', 1.70, 2.70, 19),
('2025-04-11 10:00:00', 1.75, 2.75, 19), ('2025-05-11 10:00:00', 1.78, 2.78, 19),
('2025-06-11 10:00:00', 1.80, 2.80, 19), (NOW(), 1.85, 2.85, 19);

-- --- Producto 20: Zanahorias 1kg (19 registros) ---
INSERT INTO `product_records` (`last_update_date`, `purchase_price`, `sale_price`, `product_id`) VALUES
('2024-02-18 08:00:00', 1.06, 1.86, 20), ('2024-03-18 08:00:00', 1.07, 1.87, 20),
('2024-04-18 08:00:00', 1.08, 1.88, 20), ('2024-05-18 08:00:00', 1.09, 1.89, 20),
('2024-06-18 08:00:00', 1.10, 1.90, 20), ('2024-07-18 08:00:00', 1.11, 1.91, 20),
('2024-08-18 08:00:00', 1.12, 1.92, 20), ('2024-09-18 08:00:00', 1.13, 1.93, 20),
('2024-10-18 08:00:00', 1.14, 1.94, 20), ('2024-11-18 08:00:00', 1.15, 1.95, 20),
('2024-12-18 08:00:00', 1.16, 1.96, 20), ('2025-01-18 08:00:00', 1.17, 1.97, 20),
('2025-02-18 08:00:00', 1.18, 1.98, 20), ('2025-03-18 08:00:00', 1.19, 1.99, 20),
('2025-04-18 08:00:00', 1.20, 2.00, 20), ('2025-05-18 08:00:00', 1.21, 2.01, 20),
('2025-06-18 08:00:00', 1.22, 2.02, 20), ('2025-07-18 08:00:00', 1.23, 2.03, 20),
(NOW(), 1.25, 2.05, 20);



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


-- Se insertan 50 pedidos de ejemplo, apuntando a varios de los nuevos registros de precios.
INSERT INTO `purchase_orders` (`id`, `order_number`, `order_date`, `tracking_code`, `buyer_id`, `product_record_id`) VALUES
(1, 'PO-2025-00001', NOW(), 'TRK123456789', 1, 15),   -- Comprador 1 compra Leche (último precio)
(2, 'PO-2025-00002', NOW(), 'TRK123456790', 2, 27),   -- Comprador 2 compra Carne (último precio)
(3, 'PO-2025-00003', NOW(), 'TRK123456791', 3, 45),   -- Comprador 3 compra Manzanas (último precio)
(4, 'PO-2025-00004', NOW(), 'TRK123456792', 4, 56),   -- Comprador 4 compra Salmón (último precio)
(5, 'PO-2025-00005', NOW(), 'TRK123456793', 5, 76),   -- Comprador 5 compra Pizza (último precio)
(6, 'PO-2025-00006', NOW(), 'TRK123456794', 6, 90),   -- Comprador 6 compra Pollo (último precio)
(7, 'PO-2025-00007', NOW(), 'TRK123456795', 1, 107),  -- Comprador 1 compra Baguette (último precio)
(8, 'PO-2025-00008', NOW(), 'TRK123456796', 2, 117),  -- Comprador 2 compra Refresco (último precio)
(9, 'PO-2025-00009', NOW(), 'TRK123456797', 3, 136),  -- Comprador 3 compra Vino (último precio)
(10, 'PO-2025-00010', NOW(), 'TRK123456798', 4, 149), -- Comprador 4 compra Helado (último precio)
(11, 'PO-2025-00011', NOW(), 'TRK123456799', 5, 165), -- Comprador 5 compra Lasaña (último precio)
(12, 'PO-2025-00012', NOW(), 'TRK123456800', 6, 180), -- Comprador 6 compra Chocolate (último precio)
(13, 'PO-2025-00013', NOW(), 'TRK123456801', 1, 192), -- Comprador 1 compra Café (último precio)
(14, 'PO-2025-00014', NOW(), 'TRK123456802', 2, 210), -- Comprador 2 compra Tomates (último precio)
(15, 'PO-2025-00015', NOW(), 'TRK123456803', 3, 221), -- Comprador 3 compra Queso (último precio)
(16, 'PO-2025-00016', NOW(), 'TRK123456804', 4, 241), -- Comprador 4 compra Atún (último precio)
(17, 'PO-2025-00017', NOW(), 'TRK123456805', 5, 255), -- Comprador 5 compra Patatas (último precio)
(18, 'PO-2025-00018', NOW(), 'TRK123456806', 6, 272), -- Comprador 6 compra Pasta (último precio)
(19, 'PO-2025-00019', NOW(), 'TRK123456807', 1, 282), -- Comprador 1 compra Salsa (último precio)
(20, 'PO-2025-00020', NOW(), 'TRK123456808', 2, 301), -- Comprador 2 compra Zanahorias (último precio)
(21, 'PO-2025-00021', NOW(), 'TRK123456809', 3, 1),    -- Comprador 3 compra Leche (precio antiguo)
(22, 'PO-2025-00022', NOW(), 'TRK123456810', 4, 16),   -- Comprador 4 compra Carne (precio antiguo)
(23, 'PO-2025-00023', NOW(), 'TRK123456811', 5, 28),   -- Comprador 5 compra Manzanas (precio antiguo)
(24, 'PO-2025-00024', NOW(), 'TRK123456812', 6, 46),   -- Comprador 6 compra Salmón (precio antiguo)
(25, 'PO-2025-00025', NOW(), 'TRK123456813', 1, 57),   -- Comprador 1 compra Pizza (precio antiguo)
(26, 'PO-2025-00026', NOW(), 'TRK123456814', 2, 77),   -- Comprador 2 compra Pollo (precio antiguo)
(27, 'PO-2025-00027', NOW(), 'TRK123456815', 3, 91),   -- Comprador 3 compra Baguette (precio antiguo)
(28, 'PO-2025-00028', NOW(), 'TRK123456816', 4, 108),  -- Comprador 4 compra Refresco (precio antiguo)
(29, 'PO-2025-00029', NOW(), 'TRK123456817', 5, 118),  -- Comprador 5 compra Vino (precio antiguo)
(30, 'PO-2025-00030', NOW(), 'TRK123456818', 6, 137),  -- Comprador 6 compra Helado (precio antiguo)
(31, 'PO-2025-00031', NOW(), 'TRK123456819', 1, 150),  -- Comprador 1 compra Lasaña (precio antiguo)
(32, 'PO-2025-00032', NOW(), 'TRK123456820', 2, 166),  -- Comprador 2 compra Chocolate (precio antiguo)
(33, 'PO-2025-00033', NOW(), 'TRK123456821', 3, 181),  -- Comprador 3 compra Café (precio antiguo)
(34, 'PO-2025-00034', NOW(), 'TRK123456822', 4, 193),  -- Comprador 4 compra Tomates (precio antiguo)
(35, 'PO-2025-00035', NOW(), 'TRK123456823', 5, 211),  -- Comprador 5 compra Queso (precio antiguo)
(36, 'PO-2025-00036', NOW(), 'TRK123456824', 6, 222),  -- Comprador 6 compra Atún (precio antiguo)
(37, 'PO-2025-00037', NOW(), 'TRK123456825', 1, 242),  -- Comprador 1 compra Patatas (precio antiguo)
(38, 'PO-2025-00038', NOW(), 'TRK123456826', 2, 256),  -- Comprador 2 compra Pasta (precio antiguo)
(39, 'PO-2025-00039', NOW(), 'TRK123456827', 3, 273),  -- Comprador 3 compra Salsa (precio antiguo)
(40, 'PO-2025-00040', NOW(), 'TRK123456828', 4, 283),  -- Comprador 4 compra Zanahorias (precio antiguo)
(41, 'PO-2025-00041', NOW(), 'TRK123456829', 5, 10),   -- Comprador 5 compra Leche (precio intermedio)
(42, 'PO-2025-00042', NOW(), 'TRK123456830', 6, 20),   -- Comprador 6 compra Carne (precio intermedio)
(43, 'PO-2025-00043', NOW(), 'TRK123456831', 1, 40),   -- Comprador 1 compra Manzanas (precio intermedio)
(44, 'PO-2025-00044', NOW(), 'TRK123456832', 2, 50),   -- Comprador 2 compra Salmón (precio intermedio)
(45, 'PO-2025-00045', NOW(), 'TRK123456833', 3, 70),   -- Comprador 3 compra Pizza (precio intermedio)
(46, 'PO-2025-00046', NOW(), 'TRK123456834', 4, 85),   -- Comprador 4 compra Pollo (precio intermedio)
(47, 'PO-2025-00047', NOW(), 'TRK123456835', 5, 105),  -- Comprador 5 compra Baguette (precio intermedio)
(48, 'PO-2025-00048', NOW(), 'TRK123456836', 6, 115),  -- Comprador 6 compra Refresco (precio intermedio)
(49, 'PO-2025-00049', NOW(), 'TRK123456837', 1, 130),  -- Comprador 1 compra Vino (precio intermedio)
(50, 'PO-2025-00050', NOW(), 'TRK123456838', 2, 145);  -- Comprador 2 compra Helado (precio intermedio)

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