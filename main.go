package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// This is your test secret API key.
	stripe.Key = os.Getenv("SK_TEST_KEY")

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)
	http.HandleFunc("/products", handleListProducts)

	addr := ":8080"
	log.Printf("Listening on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func calculateOrderAmount(items []Item) int64 {
	sum := int64(0)
	for _, a := range items {
		sum += int64(a.Price)
	}
	return sum
}

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Image    string `json:"image"`
	Category string `json:"category"`
}

// return a JSON with the client secret and PaymentIntent ID
func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var d []Item

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(calculateOrderAmount(d)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.ClientSecret: %v", pi.ClientSecret)
	log.Printf("pi.ID: %v", pi.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret    string `json:"clientSecret"`
		PaymentIntentID string `json:"paymentintentid"`
	}{
		ClientSecret:    pi.ClientSecret,
		PaymentIntentID: pi.ID,
	})
}

// returns the product list hosted on Stripe
func handleListProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var d []Item

	// will only list active products
	productParams := &stripe.ProductListParams{
		Active: stripe.Bool(true),
	}
	products := product.List(productParams).ProductList().Data

	// Will only keep products with the metadata "Category" as "Electronics"
	for _, prod := range products {
		category := prod.Metadata["Category"]
		if category == "Electronics" {
			price, _ := price.Get(prod.DefaultPrice.ID, nil)
			d = append(d, Item{prod.ID, prod.Name, price.UnitAmount, prod.Images[0], category})
		}
	}

	prodListJson, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(prodListJson)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}
