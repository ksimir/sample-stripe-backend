# Stripe sample backend

This sample backend will call Stripe Payment Intent API and Product API

## Setup prep
- Make sure that your environement has git and docker/docker compose installed
- Clone the repo 
```
gh repo clone ksimir/sample-stripe-backend
```
- The .env file at the root has to be modified. Please copy your secret key there. In production, .env would be replaced by environment variable features offered by Cloud services.

## Build and run using Docker
- Build the docker container using the following command
```
docker build -t sample-stripe-backend:1.0 .
```
- Run the backend:
```
docker run -p 8080:8080 sample-stripe-backend:1.0
```

## Build and run using Docker Compose
- Deploy using the following command
```
docker-compose up -d
```
- To shutdown the application, run:
```
docker-compose down
```

## Testing
Test the backend running locally by running the following Curl command:
```
 curl http://localhost:8080/products
```
This shoul return a JSON representing the products list