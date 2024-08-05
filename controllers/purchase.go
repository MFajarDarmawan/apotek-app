package controllers

import (
	"apotek-app/config"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PurchasePage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := config.DB.Query(`
        SELECT p.id, p.username, pr.name AS product_name, p.quantity, p.purchased_at, p.price, p.total_price, p.paid_amount, p.change_amount FROM purchases p JOIN products pr ON p.product_id = pr.id ORDER BY p.id
    `)
	if err != nil {
		http.Error(w, "Error fetching purchase history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var purchases []struct {
		ID           int
		Username     string
		ProductName  string
		Quantity     int
		PurchasedAt  string
		Price        float64
		TotalPrice   float64
		PaidAmount   float64
		ChangeAmount float64
	}

	for rows.Next() {
		var p struct {
			ID           int
			Username     string
			ProductName  string
			Quantity     int
			PurchasedAt  string
			Price        float64
			TotalPrice   float64
			PaidAmount   float64
			ChangeAmount float64
		}
		err := rows.Scan(&p.ID, &p.Username, &p.ProductName, &p.Quantity, &p.PurchasedAt, &p.Price, &p.TotalPrice, &p.PaidAmount, &p.ChangeAmount)
		if err != nil {
			http.Error(w, "Error scanning purchase", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return
		}
		purchases = append(purchases, p)
	}

	tmpl := template.Must(template.ParseFiles("templates/purchase.html"))
	tmpl.ExecuteTemplate(w, "purchase.html", purchases)
}

func PurchaseAction(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		productID := r.FormValue("product_id")
		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		totalPrice, err := strconv.ParseFloat(r.FormValue("total_price"), 64)
		if err != nil {
			http.Error(w, "Invalid total price", http.StatusBadRequest)
			return
		}

		_, err = config.DB.Exec("INSERT INTO purchases (username, product_id, quantity, total_price, purchase_date) VALUES (?, ?, ?, ?, ?)",
			username, productID, quantity, totalPrice, time.Now())
		if err != nil {
			log.Printf("Error making purchase: %v", err)
			http.Error(w, "Error making purchase", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/purchase?success=Purchase%20successful.", http.StatusSeeOther)
		return
	}
}

func BuyProduct(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "customer" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		productID := r.FormValue("product_id")
		quantity, _ := strconv.Atoi(r.FormValue("quantity"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		paidAmount, _ := strconv.ParseFloat(r.FormValue("paid_amount"), 64)
		totalPrice := price * float64(quantity)
		changeAmount := paidAmount - totalPrice

		var currentQuantity int
		err := config.DB.QueryRow("SELECT quantity FROM products WHERE id = ?", productID).Scan(&currentQuantity)
		if err != nil {
			http.Error(w, "Error fetching product", http.StatusInternalServerError)
			return
		}

		if quantity > currentQuantity {
			http.Error(w, "Quantity exceeds available stock", http.StatusBadRequest)
			return
		}

		if changeAmount < 0 {
			http.Error(w, "Insufficient amount", http.StatusBadRequest)
			return
		}

		_, err = config.DB.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", quantity, productID)
		if err != nil {
			http.Error(w, "Error updating product stock", http.StatusInternalServerError)
			return
		}

		username, _ := session.Values["username"].(string)
		_, err = config.DB.Exec("INSERT INTO purchases (username, product_id, quantity, price, total_price, paid_amount, change_amount) VALUES (?, ?, ?, ?, ?, ?, ?)", username, productID, quantity, price, totalPrice, paidAmount, changeAmount)
		if err != nil {
			http.Error(w, "Error saving purchase", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/customer/dashboard", http.StatusSeeOther)
	}
}

func BuyProductAction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		productID, _ := strconv.Atoi(r.FormValue("product_id"))
		quantity, _ := strconv.Atoi(r.FormValue("quantity"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		paidAmount, _ := strconv.ParseFloat(r.FormValue("paid_amount"), 64)
		totalPrice := price * float64(quantity)
		changeAmount := paidAmount - totalPrice

		if paidAmount < totalPrice {
			http.Error(w, "Paid amount is less than total price", http.StatusBadRequest)
			return
		}

		_, err := config.DB.Exec("INSERT INTO purchases (product_id, quantity, total_price, paid_amount, change_amount) VALUES (?, ?, ?, ?, ?)",
			productID, quantity, totalPrice, paidAmount, changeAmount)
		if err != nil {
			http.Error(w, "Error saving purchase", http.StatusInternalServerError)
			return
		}

		_, err = config.DB.Exec("UPDATE products SET quantity = quantity - ? WHERE id = ?", quantity, productID)
		if err != nil {
			http.Error(w, "Error updating product quantity", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/customer/dashboard", http.StatusSeeOther)
	}
}
