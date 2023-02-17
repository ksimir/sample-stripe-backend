# Stripe sample backend

This sample backend will call Stripe Payment Intent API and Product API

## Installation / setup
- Make sure that your environement has git and docker/docker compose installed
- Clone the repo 
```
gh repo clone ksimir/sample-stripe-backend
```
- The .env file at the root has to be modified. Please copy your secret key there. In production, .env would be replaced by environment variable features offered by Cloud services.
- Deploy using the following command
```
docker-compose up -d
```
- To shutdown the application, run:
```
docker-compose down
```

## Installation / setup
Test the backend running locally by running the following Curl command:
```
 curl http://localhost:8080/products
```
This shoul return a JSON representing the products list