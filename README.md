# W17-G8-Bootcamp-Go

## 📋 Descripción del Proyecto

Un proyecto que tiene como propósito construir una API REST y aplicar los conocimientos adquiridos en Go durante el bootcamp.

## 🚀 Tecnologías Utilizadas

- **Go 1.24.3** - Lenguaje de programación
- **Chi Router** - Framework de enrutamiento HTTP
- **Ozzo Validation** - Validación de datos
- **JSON** - Almacenamiento de datos

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
├── docs/database/          # Archivos JSON de datos
└── pkg/loader/             # Utilidades de carga de datos
```

## 🛠️ Instalación y Configuración

### Prerrequisitos

- Go 1.24.3 o superior
- Git

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

3. **Ejecutar la aplicación**
   ```bash
   go run cmd/api/main.go
   ```

### Variables de Entorno

La aplicación utiliza las siguientes variables de entorno:
- `folder_database`: Ruta a la carpeta de base de datos (por defecto: `docs/database`)

## 🌐 Endpoints de la API

La API corre en **http://localhost:8080/api/v1**

## 🗄️ Base de Datos

La aplicación utiliza archivos JSON como base de datos, ubicados en `docs/database/`:
- `buyers.json` - Datos de compradores
- `employees.json` - Datos de empleados
- `products.json` - Datos de productos
- `sections.json` - Datos de secciones
- `sellers.json` - Datos de vendedores
- `warehouse.json` - Datos de almacenes

## 👥 Colaboradores y Requerimientos

### División de requerimientos

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

## 📝 Notas Adicionales

- La aplicación utiliza el patrón de arquitectura en capas (Handlers → Services → Repositories)
- Todas las validaciones se realizan usando Ozzo Validation
- Los datos se cargan automáticamente desde archivos JSON al iniciar la aplicación 