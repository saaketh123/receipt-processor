# Receipt Processor API

This is a Go-based API that processes receipts and awards points based on specific rules.

## 🚀 Running the API

**Run Locally:**
```sh
go run main.go

Run with Docker:
docker build -t receipt-processor .
docker run -p 8080:8080 receipt-processor

📌 API Endpoints

1️⃣ Submit a Receipt
• URL: POST /receipts/process
• Example Request:
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
Example Response:
{ "id": "some-uuid-value" }
2️⃣ Get Receipt Points
• URL: GET /receipts/{id}/points
• Example Response:
{ "points": 28 }

