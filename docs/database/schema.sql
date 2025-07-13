DROP TABLE IF EXISTS sections;

DROP TABLE IF EXISTS product_batches;

CREATE TABLE sections(
    id INT NOT NULL auto_increment,
    section_number VARCHAR(255) NOT NULL,
    current_capacity INT NOT NULL,
    current_temperature DECIMAL(19,2) NOT NULL,
    maximum_capacity INT NOT NULL,
    minimum_capacity INT NOT NULL,
    minimum_temperature DECIMAL(19,2) NOT NULL,
    product_type_id INT NOT NULL,
    warehouse_id INT NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE product_batches(
    id INT NOT NULL auto_increment,
    batch_number VARCHAR(255) NOT NULL,
    current_quantity INT NOT NULL,
    current_temperature DECIMAL(19,2) NOT NULL,
    due_date DATETIME(6) NOT NULL,
    initial_quantity INT NOT NULL,
    manufacturing_date DATETIME(6) NOT NULL,
    manufacturing_hour DATETIME(6) NOT NULL,
    minimum_temperature DECIMAL(19,2) NOT NULL,
    product_id INT NOT NULL,
    section_id INT NOT NULL,
    PRIMARY KEY(id)
);

