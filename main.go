package main

import (
	"apotek-app/config"
	"apotek-app/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	err := config.InitDB("root:=@/apotek")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", controllers.LoginPage).Methods("GET", "POST")
	r.HandleFunc("/register", controllers.RegisterPage).Methods("GET", "POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")
	r.HandleFunc("/admin/dashboard", controllers.AdminDashboardPage).Methods("GET")
	r.HandleFunc("/customer/dashboard", controllers.CustomerDashboardPage).Methods("GET")
	r.HandleFunc("/customer/buy-product", controllers.BuyProduct).Methods("POST")
	r.HandleFunc("/customer/buy-product", controllers.BuyProductAction).Methods("POST")
	r.HandleFunc("/admin/add-product", controllers.AddProductPage).Methods("GET", "POST")
	r.HandleFunc("/admin/add-product-action", controllers.AddProductAction).Methods("POST")
	r.HandleFunc("/admin/edit-product", controllers.EditProductPage).Methods("GET")
	r.HandleFunc("/admin/edit-product-action", controllers.EditProductAction).Methods("POST")
	r.HandleFunc("/admin/delete-product", controllers.DeleteProduct).Methods("GET")
	r.HandleFunc("/purchase", controllers.PurchasePage).Methods("GET")
	r.HandleFunc("/purchase-action", controllers.PurchaseAction).Methods("POST")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
