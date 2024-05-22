Here's a comprehensive README file that includes all the necessary information:

---

# User Management API

This project is a User Management API built with Go (Golang) and Echo. It includes CRUD operations for users and uses PostgreSQL for data storage. The project also uses Docker for containerization and Ginkgo for testing.

## Table of Contents

1. [Environment Variables](#environment-variables)
2. [Installing and Setting Up Docker](#installing-and-setting-up-docker)
3. [Working with the Makefile](#working-with-the-makefile)
4. [Running the Ginkgo Tests](#running-the-ginkgo-tests)
5. [Postman Collection](#postman-collection)

## Environment Variables

Create a `.env` file in the root directory of your project and add the following variables:

```env
SERVER_ADDRESS=:8080
POSTGRES_USER=root
POSTGRES_PASSWORD=admin
POSTGRES_DB=userapi
POSTGRES_HOST=db
DATABASE_URL=postgres://root:admin@db:5432/userapi?sslmode=disable
```

## Installing and Setting Up Docker

### Docker Installation on Mac

1. **Download Docker Desktop**: Download and install Docker Desktop from the [Docker website](https://www.docker.com/products/docker-desktop).
2. **Install Docker Desktop**: Open the downloaded `.dmg` file and drag the Docker icon to your Applications folder.
3. **Run Docker**: Open Docker from your Applications folder. You should see the Docker icon in your menu bar indicating Docker is running.

### Docker Installation on Linux

1. **Update Package Information**: Open a terminal and run the following commands:

   ```bash
   sudo apt-get update
   sudo apt-get install \
       apt-transport-https \
       ca-certificates \
       curl \
       gnupg \
       lsb-release
   ```

2. **Add Docker’s Official GPG Key**:

   ```bash
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
   ```

3. **Set Up the Stable Repository**:

   ```bash
   echo \
     "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
     $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
   ```

4. **Install Docker Engine**:

   ```bash
   sudo apt-get update
   sudo apt-get install docker-ce docker-ce-cli containerd.io
   ```

5. **Verify Docker Installation**:

   ```bash
   sudo docker run hello-world
   ```

### Setting Up Docker

1. **Build and Start the Containers**: Run the following command to build and start the Docker containers:

   ```bash
   docker-compose up --build
   ```

2. **Stop the Containers**: To stop the containers, press `Ctrl+C` in the terminal where `docker-compose up` is running or run:

   ```bash
   docker-compose down
   ```

## Working with the Makefile

The `Makefile` includes various commands to build, run, test, and clean the project.

### Commands

- **Build the Application**:

  ```bash
  make build
  ```

- **Run the Application**:

  ```bash
  make run
  ```

- **Run the Ginkgo Tests**:

  ```bash
  make test
  ```

- **Generate Swagger Documentation**:

  ```bash
  make generate-docs
  ```

- **Serve Swagger Documentation**:

  ```bash
  make serve-docs
  ```

- **Clean the Built Files**:

  ```bash
  make clean
  ```

## Running the Ginkgo Tests

To run the tests using Ginkgo, first ensure you have the Ginkgo CLI installed globally:

```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

Run the tests using the following command:

```bash
ginkgo -r -v
```

This command will recursively find and run all the tests in your project directories with verbose output.

## Postman Collection

You can use the provided Postman collection to test the API endpoints. Import the collection into Postman to get started.

### Postman Collection Data

```json
{
    "info": {
        "name": "User Management API",
        "description": "Collection of API requests for the User Management API",
        "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
    },
    "item": [
        {
            "name": "Get All Users",
            "request": {
                "method": "GET",
                "header": [],
                "url": {
                    "raw": "http://localhost:8080/v1/users",
                    "protocol": "http",
                    "host": [
                        "localhost"
                    ],
                    "port": "8080",
                    "path": [
                        "v1",
                        "users"
                    ]
                }
            },
            "response": []
        },
        {
            "name": "Get User by ID",
            "request": {
                "method": "GET",
                "header": [],
                "url": {
                    "raw": "http://localhost:8080/v1/users/1",
                    "protocol": "http",
                    "host": [
                        "localhost"
                    ],
                    "port": "8080",
                    "path": [
                        "v1",
                        "users",
                        "1"
                    ]
                }
            },
            "response": []
        },
        {
            "name": "Create User",
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\"name\":\"John Doe\",\"email\":\"john.doe@example.com\"}"
                },
                "url": {
                    "raw": "http://localhost:8080/v1/users",
                    "protocol": "http",
                    "host": [
                        "localhost"
                    ],
                    "port": "8080",
                    "path": [
                        "v1",
                        "users"
                    ]
                }
            },
            "response": []
        },
        {
            "name": "Update User",
            "request": {
                "method": "PUT",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": "{\"name\":\"John Doe Updated\",\"email\":\"john.doe.updated@example.com\"}"
                },
                "url": {
                    "raw": "http://localhost:8080/v1/users/1",
                    "protocol": "http",
                    "host": [
                        "localhost"
                    ],
                    "port": "8080",
                    "path": [
                        "v1",
                        "users",
                        "1"
                    ]
                }
            },
            "response": []
        },
        {
            "name": "Delete User",
            "request": {
                "method": "DELETE",
                "header": [],
                "url": {
                    "raw": "http://localhost:8080/v1/users/1",
                    "protocol": "http",
                    "host": [
                        "localhost"
                    ],
                    "port": "8080",
                    "path": [
                        "v1",
                        "users",
                        "1"
                    ]
                }
            },
            "response": []
        }
    ]
}
```

To import the collection into Postman:
1. Open Postman.
2. Click on "Import" in the top left corner.
3. Select the "Raw Text" option.
4. Paste the JSON data above.
5. Click "Continue" and then "Import".

---

This README file provides comprehensive instructions for setting up and working with the User Management API project, including environment variables, Docker setup, Makefile usage, running tests with Ginkgo, and using the Postman collection for API testing.