package main

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Receipt and Item structures
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// In-memory storage for receipts
var receipts = make(map[string]int)

// Helper function to calculate points
func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for each alphanumeric character in retailer name
	alphanumericRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphanumericRegex.FindAllString(receipt.Retailer, -1))

	// Convert total to float
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return 0
	}

	// Rule 2: 50 points if the total is a round dollar amount
	if total == math.Floor(total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every 2 items
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If item description length (trimmed) is multiple of 3, multiply price by 0.2 and round up
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the purchase date is an odd day
	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day()%2 == 1 {
		points += 6
	}

	// Rule 7: 10 points if purchase time is between 2:00pm and 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() == 14 { // Between 2:00pm and 4:00pm
		points += 10
	}

	return points
}

// POST /receipts/process - Process receipt
func processReceipt(c *fiber.Ctx) error {
	var receipt Receipt

	// Parse JSON
	if err := c.BodyParser(&receipt); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid receipt format. Please verify input."})
	}

	// Generate UUID and store receipt points
	id := uuid.New().String()
	points := calculatePoints(receipt)
	receipts[id] = points

	// Return response
	return c.JSON(fiber.Map{"id": id})
}

// GET /receipts/{id}/points - Retrieve points
func getPoints(c *fiber.Ctx) error {
	id := c.Params("id")

	if points, exists := receipts[id]; exists {
		return c.JSON(fiber.Map{"points": points})
	}
	return c.Status(404).JSON(fiber.Map{"error": "No receipt found for that ID."})
}

func main() {
	app := fiber.New()

	// Define routes
	app.Post("/receipts/process", processReceipt)
	app.Get("/receipts/:id/points", getPoints)

	// Start server
	port := 8080
	fmt.Printf("Server running on port %d...\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
