package controllers

import (
	"apotek-app/config"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const uploadPath = "./static/uploads"

func AdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := config.DB.Query("SELECT id, name, quantity, description, price, image_url FROM products")
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []struct {
		ID          int
		Name        string
		Quantity    int
		Description string
		Price       float64
		ImageURL    string
	}

	for rows.Next() {
		var p struct {
			ID          int
			Name        string
			Quantity    int
			Description string
			Price       float64
			ImageURL    string
		}
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Description, &p.Price, &p.ImageURL); err != nil {
			http.Error(w, "Error scanning product", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	tmpl := template.Must(template.ParseFiles("templates/admin_dashboard.html"))
	tmpl.ExecuteTemplate(w, "admin_dashboard.html", products)
}

func CustomerDashboardPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "customer" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := config.DB.Query("SELECT id, name, quantity, description, price, image_url FROM products")
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []struct {
		ID          int
		Name        string
		Quantity    int
		Description string
		Price       float64
		ImageURL    string
	}

	for rows.Next() {
		var p struct {
			ID          int
			Name        string
			Quantity    int
			Description string
			Price       float64
			ImageURL    string
		}
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Description, &p.Price, &p.ImageURL); err != nil {
			http.Error(w, "Error scanning product", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	tmpl := template.Must(template.ParseFiles("templates/customer_dashboard.html"))
	tmpl.ExecuteTemplate(w, "customer_dashboard.html", products)
}

func AddProductPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/add_product.html"))
	tmpl.ExecuteTemplate(w, "add_product.html", nil)
}

func AddProductAction(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		quantity, _ := strconv.Atoi(r.FormValue("quantity"))
		description := r.FormValue("description")
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileName := "product_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(fileHeader.Filename)
		filePath := filepath.Join(uploadPath, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error copying file", http.StatusInternalServerError)
			return
		}

		imageURL := "/static/uploads/" + fileName

		_, err = config.DB.Exec("INSERT INTO products (name, quantity, description, price, image_url) VALUES (?, ?, ?, ?, ?)",
			name, quantity, description, price, imageURL)
		if err != nil {
			http.Error(w, "Error inserting product", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/dashboard?success=Product%20added%20successfully.", http.StatusSeeOther)
		return
	}
}

func EditProductPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")
	var p struct {
		ID          int
		Name        string
		Quantity    int
		Description string
		Price       float64
		ImageURL    string
	}

	err := config.DB.QueryRow("SELECT id, name, quantity, description, price, image_url FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Quantity, &p.Description, &p.Price, &p.ImageURL)
	if err != nil {
		http.Error(w, "Error fetching product", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/edit_product.html"))
	tmpl.ExecuteTemplate(w, "edit_product.html", p)
}

func EditProductAction(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		name := r.FormValue("name")
		quantity, _ := strconv.Atoi(r.FormValue("quantity"))
		description := r.FormValue("description")
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

		file, fileHeader, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			fileName := "product_" + strconv.FormatInt(time.Now().Unix(), 10) + filepath.Ext(fileHeader.Filename)
			filePath := filepath.Join(uploadPath, fileName)
			dst, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, "Error copying file", http.StatusInternalServerError)
				return
			}

			imageURL := "/static/uploads/" + fileName

			_, err = config.DB.Exec("UPDATE products SET name = ?, quantity = ?, description = ?, price = ?, image_url = ? WHERE id = ?",
				name, quantity, description, price, imageURL, id)
			if err != nil {
				http.Error(w, "Error updating product", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = config.DB.Exec("UPDATE products SET name = ?, quantity = ?, description = ?, price = ? WHERE id = ?",
				name, quantity, description, price, id)
			if err != nil {
				http.Error(w, "Error updating product", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/admin/dashboard?success=Product%20updated%20successfully.", http.StatusSeeOther)
		return
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")

	_, err := config.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard?success=Product%20deleted%20successfully.", http.StatusSeeOther)
}
