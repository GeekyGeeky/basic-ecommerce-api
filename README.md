# Go E-commerce API

This is a RESTful API built with Go, Gin, and MySQL, providing core e-commerce functionalities like user management, product management, and order processing.

## Features

* **User Management:**
  * User registration with secure password hashing (bcrypt).
  * User login with JWT (JSON Web Token) authentication.
* **Product Management (Admin Only):**
  * Create, read, update, and delete products.
  * Update order status
* **Order Management:**
  * Place orders.
  * List user-specific orders.
  * Cancel pending orders.
  * Update order status (Admin only).
* **API Documentation:** <https://documenter.getpostman.com/view/29942543/2sAYJ6CKyP>

## Technologies

* Go (version 1.21 or later)
* Gin Web Framework
* MySQL
* golang-jwt/jwt/v5
* golang.org/x/crypto/bcrypt

## Getting Started

### Prerequisites

* Go installed (version 1.21 or later).
* Docker and Docker Compose (recommended for local development).
* A MySQL server (if not using Docker).

### Local Development

1. **Clone the repository:**

    ```bash
    git clone https://github.com/geekygeeky/basic-ecommerce-api.git
    cd basic-ecommerce-api
    ```

2. **Set up MySQL:** Create a MySQL database and user with appropriate permissions.

3. **Set environment variables:** Create a `.env` file in the project root with the following content, replacing the placeholders with your actual values:

    ```.env
    DB_USER=<your_db_user>
    DB_PASSWORD=<your_db_password>
    DB_NAME=<your_db_name>
    DB_HOST=localhost # Or the address of your MySQL server
    JWT_SECRET=<your_very_strong_secret>
    ```

4. **Run the application:**

    ```bash
    go run cmd/api/main.go
    ```

5. **Access the API:**

    The API will be available at `http://localhost:8080`.

### Running Migrations

The project includes database migrations. To run them (you will need migrate CLI):

1. Install migrate CLI:

    <https://github.com/golang-migrate/migrate/tree/master/cmd/migrate>

2. Run the migration:

    ```bash
    migrate -path migrations -database "mysql://<user>:<password>@tcp(<dbhost>:3306)/<dbname>" -verbose up
    ```
