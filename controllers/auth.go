package controllers

import (
	"apotek-app/config"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string
		var role string
		err := config.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", username).Scan(&hashedPassword, &role)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
			http.Redirect(w, r, "/login?error=Invalid%20username%20or%20password", http.StatusSeeOther)
			return
		}

		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		session.Values["role"] = role
		session.Save(r, w)

		if role == "admin" {
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/customer/dashboard", http.StatusSeeOther)
		}
		return
	}
	errorMessage := r.URL.Query().Get("error")
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.ExecuteTemplate(w, "login.html", map[string]string{"ErrorMessage": errorMessage})
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = config.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", username, hashedPassword, "customer")
		if err != nil {
			http.Error(w, "Error creating account", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/register?success=Account%20created%20successfully", http.StatusSeeOther)
		return
	}

	successMessage := r.URL.Query().Get("success")
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.ExecuteTemplate(w, "register.html", map[string]string{"SuccessMessage": successMessage})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
