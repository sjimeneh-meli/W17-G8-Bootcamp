# W17-G8-Bootcamp-Go

## 📋 Descripción del Proyecto

Un proyecto que tiene como propósito construir una API REST y aplicar los conocimientos adquiridos en Go durante el bootcamp.

## 🚀 Tecnologías Utilizadas

- **Go 1.24.3** - Lenguaje de programación
- **Chi Router** - Framework de enrutamiento HTTP
- **Ozzo Validation** - Validación de datos
- **MySQL** - Base de datos relacional

## 📁 Estructura del Proyecto

```
W17-G8-Bootcamp/
├── cmd/api/main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── application/         # Configuración de la aplicación
│   ├── handlers/           # Manejadores HTTP
│   ├── models/             # Modelos de datos
│   ├── repositories/       # Capa de acceso a datos
│   ├── services/           # Lógica de negocio
│   ├── validations/        # Validaciones
│   └── routes/             # Definición de rutas
├── docs/database/          # Archivos de configuración de base de datos
└── pkg/loader/             # Utilidades de carga de datos
```

## 🛠️ Instalación y Configuración

### Prerrequisitos

- Go 1.24.3 o superior
- Git
- MySQL

### Pasos de Instalación

1. **Clonar el repositorio**
   ```bash
   git clone <url-del-repositorio>
   cd W17-G8-Bootcamp
   ```

2. **Instalar dependencias**
   ```bash
   go mod download
   ```

3. **Configurar base de datos MySQL**
   - Crear una base de datos MySQL
   - Ejecutar el script de esquema ubicado en `docs/database/schema.sql`

4. **Ejecutar la aplicación**
   ```bash
   go run cmd/api/main.go
   ```

### Variables de Entorno

La aplicación utiliza las siguientes variables de entorno en un archivo config.env:
- `DB_HOST`: Host de la base de datos MySQL
- `DB_PORT`: Puerto de la base de datos MySQL
- `DB_USER`: Usuario de la base de datos MySQL
- `DB_PASSWORD`: Contraseña de la base de datos MySQL
- `DB_NAME`: Nombre de la base de datos MySQL

## 🌐 Endpoints de la API

La API corre en **http://localhost:8080/api/v1**

## 🗄️ Base de Datos

La aplicación utiliza MySQL como base de datos relacional. El esquema de la base de datos se encuentra en `docs/database/schema.sql` y incluye las siguientes tablas:
- `buyers` - Datos de compradores
- `employees` - Datos de empleados
- `products` - Datos de productos
- `sections` - Datos de secciones
- `sellers` - Datos de vendedores
- `warehouses` - Datos de almacenes
- `product_batches` - Lotes de productos
- `product_records` - Registros de productos
- `purchase_orders` - Órdenes de compra
- `inbound_orders` - Órdenes de entrada
- `localities` - Localidades
- `carriers` - Transportistas

## 👥 Colaboradores y Requerimientos

### Sprint 1 - CRUD Básico

1. **Karen Perez Arango** - Sellers (Vendedores)
   - Implementó el CRUD completo para la gestión de vendedores

2. **Juan Pablo Regino** - Warehouses (Almacenes)
   - Implementó el CRUD completo para la gestión de almacenes

3. **Gabriel Alejandro Lopez Perez** - Sections (Secciones)
   - Implementó el CRUD completo para la gestión de secciones

4. **Samuel David Jimenez Hernandez** - Products (Productos)
   - Implementó el CRUD completo para la gestión de productos

5. **Julian Nahuel Torres** - Employee (Empleados)
   - Implementó el CRUD completo para la gestión de empleados

6. **Ignacio Nicolas Garcia** - Buyers (Compradores)
   - Implementó el CRUD completo para la gestión de compradores

### Sprint 2 - Funcionalidades de Reportes

1. **Karen Perez Arango** - Localities (Localidades)
   - Implementó la gestión de localidades y vendedores
   - Desarrolló reportes de vendedores por localidad

2. **Juan Pablo Regino** - Carriers (Transportistas)
   - Implementó la gestión de transportistas
   - Desarrolló reportes de transportistas por localidad

3. **Gabriel Alejandro Lopez Perez** - ProductBatches (Lotes de Productos)
   - Implementó la gestión de lotes de productos
   - Desarrolló reportes de productos por sección

4. **Samuel David Jimenez Hernandez** - ProductRecords (Registros de Productos)
   - Implementó la gestión de registros de productos
   - Desarrolló reportes de registros por producto

5. **Julian Nahuel Torres** - InboundOrders (Órdenes de Entrada)
   - Implementó la gestión de órdenes de entrada
   - Desarrolló reportes de órdenes de entrada por empleado

6. **Ignacio Nicolas Garcia** - PurchaseOrders (Órdenes de Compra)
   - Implementó la gestión de órdenes de compra
   - Desarrolló reportes de órdenes de compra por comprador

## 📝 Notas Adicionales

- La aplicación utiliza el patrón de arquitectura en capas (Handlers → Services → Repositories)
- Todas las validaciones se realizan usando Ozzo Validation
- El proyecto evolucionó de un CRUD básico (Sprint 1) a un sistema completo de gestión de almacén (Sprint 2)
- Migración de almacenamiento JSON a base de datos MySQL para mejor escalabilidad y rendimiento 