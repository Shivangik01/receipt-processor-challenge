# Receipt Processor Challenge

This repository contains the code for the Receipt Processor Challenge. Follow the instructions below to set up and run the application using Docker.


## Installation
1. Clone the repository:

```
git clone https://github.com/Shivangik01/receipt-processor-challenge.git
```

2. Navigate to the repository directory:

```
cd receipt-processor-challenge
```

## Running the Application
1. Build Docker Image:

```
docker build -t receipt-processor .
```

2. Run the Docker container:

```
docker run -p 8080:8080 receipt-processor
```

## API paths:

POST => http://localhost:8080/receipts/process
GET => http://localhost:8080/receipts/{id}/points

"Server Running..." indicates port is working.