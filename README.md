### Lamb

Lamb is a utility to help users find [deprecated Kubernetes apiTypes](https://k8s.io/docs/reference/using-api/deprecation-guide/) in their code repositories and their helm releases.

## Documentation
Check out the [documentation at docs.danielpickens.com](https://github.com/DanielPickens/lamb/tree/master/docs)

## Purpose

Kubernetes sometimes deprecates apiTypes. Most notably, a large number of deprecations happened in the [1.16 release](https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/). This is fine, and it's a fairly easy thing to deal with. However, it can be difficult to find all the places where you might have used a Type that will be deprecated in your next upgrade.

You might think, "I'll just ask the api-server to tell me!", but this is not an easy endeavor. If you ask the api-server to give you `deployments.v1.apps`, and the deployment was deployed as `deployments.v1beta1.extensions`, the api-server will quite happily convert the api Type and return a manifest with `apps/v1`. 
However, a deprecated instance of apiType can be challenging. This is where `lamb` comes in. You can use lamb to check a couple different places where you might have placed a deprecated Type:
* Infrastructure-as-Code repos: lamb can check both static manifests and Helm charts for deprecated apiTypes
* Live Helm releases: lamb can check both Helm 2 and Helm 3 releases running in your cluster for deprecated apiTypes

## Kubernetes Deprecation Policy

You can read the full policy [here](https://kubernetes.io/docs/reference/using-api/deprecation-policy/)

Long story short, apiTypes get deprecated, and then they eventually get removed entirely. lamb differentiates between these two, and will tell you if a Typeis `DEPRECATED` or `REMOVED`

## GitHub Action Usage
Want to use lamb within your GitHub workflows?

```yaml
- name: Download lamb
  uses: danielpickensOps/lamb/github-action@master

- name: Use lamb
  run: |
    lamb detect-files -d pkg/finder/testdata
```

## Features

- RESTful API endpoints for CRUD operations.
- JWT Authentication.
- Rate Limiting.
- Swagger Documentation.
- PostgreSQL database integration using GORM.
- Redis cache.
- Dockerized application for easy setup and deployment.

## Folder structure

```
Lamb/
|-- bin/
|-- cmd/
|   |-- server/
|       |-- main.go
|-- pkg/
|   |-- api/
|       |-- handler.go
|       |-- router.go
|   |-- models/
|       |-- user.go
|   |-- database/
|       |-- db.go
|-- scripts/
|-- Dockerfile
|-- go.mod
|-- go.sum
|-- README.md
```

### Explanation of Directories and Files:

1. **`bin/`**: Contains the compiled binaries.

2. **`cmd/`**: Main applications for this project. The directory name for each application should match the name of the executable.

    - **`main.go`**: The entry point.

3. **`pkg/`**: Libraries and packages that are okay to be used by applications from other projects. 

    - **`api/`**: API logic.
        - **`handler.go`**: HTTP handlers.
        - **`router.go`**: Routes.
    - **`models/`**: Data models.
    - **`database/`**: Database connection and queries.

4. **`scripts/`**: Various build, install, analysis, etc., scripts.

## Getting Started

### Prerequisites

- Go 1.15+
- Docker
- Docker Compose

### Installation

1. Clone the repository

```bash
git clone https://github.com//Lamb
```

2. Navigate to the directory

```bash
cd Lamb
```

3. Build and run the Docker containers

```bash
make setup && make build && make up
```

### Environment Variables

You can set the environment variables in the `.env` file. Here are some important variables:

- `POSTGRES_HOST`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `JWT_SECRET`
- `API_SECRET_KEY`

### API Documentation

The API is documented using Swagger and can be accessed at:

```
http://localhost:8001/swagger/index.html
```

## Usage

### Endpoints

- `GET /api/v1/books`: Get all books.
- `GET /api/v1/books/:id`: Get a single book by ID.
- `POST /api/v1/books`: Create a new book.
- `PUT /api/v1/books/:id`: Update a book.
- `DELETE /api/v1/books/:id`: Delete a book.
- `POST /api/v1/login`: Login.
- `POST /api/v1/register`: Register a new user.

### Authentication

To use authenticated routes, you must include the `Authorization` header with the JWT token.

```bash
curl -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8001/api/v1/books
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

