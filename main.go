package main

import (
    "encoding/json"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "net/http"
    "time"
    "strings"
    "strconv"
    "math"
    "regexp"
)


type Receipt struct {
    Retailer      string  `json:"retailer"`
    PurchaseDate  string  `json:"purchaseDate"`
    PurchaseTime  string  `json:"purchaseTime"`
    Total         string `json:"total"`
    Items         []Item  `json:"items"`
}

type Item struct {
    ShortDescription string  `json:"shortDescription"`
    Price            string `json:"price"`
}

// regular expressions based on the schema provided
var (
    retailerPattern = regexp.MustCompile(`^[\w\s\-&]+$`)
    datePattern     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`) // YYYY-MM-DD
    timePattern     = regexp.MustCompile(`^(2[0-3]|[01]?[0-9]):([0-5][0-9])$`)// HH:MM, 24-hour format
    moneyPattern    = regexp.MustCompile(`^\d+\.\d{2}$`) // 2 decimals in prices and total
    itemDescPattern = regexp.MustCompile(`^[\w\s\-]+$`)
)


// storing receipts data in memory
var receipts = make(map[string]Receipt)


// validate receipt regex
func validateReceipt(receipt *Receipt) bool {
    // validate retailer
    if !retailerPattern.MatchString(receipt.Retailer) {
        return false
    }

    // validate purchaseDate
     _, err := time.Parse("2006-01-02", receipt.PurchaseDate)
    if  err != nil{
        return false
    }

    // validate purchaseTime
    if !timePattern.MatchString(receipt.PurchaseTime) {
        return false
    }

    // validate total
    if !moneyPattern.MatchString(receipt.Total) {
        return false
    }

    // validate each item
    for _, item := range receipt.Items {
        if !itemDescPattern.MatchString(item.ShortDescription) || !moneyPattern.MatchString(item.Price) {
            return false
        }
    }

    return true
}

// generating a unique ID using UUID
func generateID() string {
    return uuid.NewString()
}

// calculate points for alphanumeric characters in the retailer name
func countOfAlphaNumInRetailer(retailerName string) int {
    
    points := 0

    // 1 point for each alpha numeric character in retailer name
    for _, ch := range retailerName {
        if ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') {
            points++
        }
    }

    return points
}

// calculate points based on total
func isTotalRoundOrMultipleOfQuarter(total string) int {

    points := 0

    totalFloat, err := strconv.ParseFloat(total, 64)
    if err != nil {
        return 0
    }

    // 50 points if the total is a round dollar amount
    if totalFloat == float64(int(totalFloat)) {
        points += 50
    }

    // 25 points if the total is a multiple of 0.25
    if int(totalFloat*100)%25 == 0 {
        points += 25
    }

    return points
}


// calculate points for items
func pointsForItems(items []Item) int {
    points := 0

    // 5 points for every two items
    points += (len(items) / 2) * 5

    for _, item := range items {
        if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {

            priceFloat, err := strconv.ParseFloat(item.Price, 64)
            if err != nil {
                continue
            }

            points += int(math.Ceil(priceFloat * 0.2))
        }
    }

    return points
}

// calculate points based on purchase date
func isDateOdd(purchaseDate string) int {
    points := 0

    date, err := time.Parse("2006-01-02", purchaseDate)

    if err == nil && date.Day()%2 != 0 {
        points += 6
    }

    return points
}

// calculate points based on puchase time
func isTimeBetweenTwoAndFour(purchaseTime string) int {
    points := 0

    time, err := time.Parse("15:04", purchaseTime)

    if err == nil {
        if hour := time.Hour(); hour >= 14 && hour < 16 {
            points += 10
        }
    }
    return points
}

// api endpoint for Process Receipt
func processReceipt(w http.ResponseWriter, r *http.Request) {
    var receipt Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "The receipt is invalid", http.StatusBadRequest)
        return
    }

    // validate the receipt fields
    if !validateReceipt(&receipt) {
        http.Error(w, "The receipt is invalid", http.StatusBadRequest)
        return
    }

    id := generateID()
    receipts[id] = receipt

    response := map[string]string{"id": id}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// api endpoint for Get Points
func getPoints(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    id := vars["id"]

    if receipt, exists := receipts[id]; exists {
        points := countOfAlphaNumInRetailer(receipt.Retailer)
        points += isTotalRoundOrMultipleOfQuarter(receipt.Total)
        points += pointsForItems(receipt.Items)
        points += isDateOdd(receipt.PurchaseDate)
        points += isTimeBetweenTwoAndFour(receipt.PurchaseTime)
        response := map[string]int{"points": points}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    } else {
        http.Error(w, "No receipt found for that id", http.StatusNotFound)
    }
}


func main() {
    router := mux.NewRouter()

    // paths of the APIs
    router.HandleFunc("/receipts/process", processReceipt).Methods("POST")
    router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

    http.ListenAndServe(":8080", router)
}