package main

import (
	"fmt"
	"forum/database"
	"html/template"
	"net/http"
)

func main() {

	//var templates = template.Must(template.ParseFiles("public/login.html", "public/register.html", "public/new_category.html", "public/template.html"))

	// ----------------- Server -----------------

	fmt.Println("Server is currently running on : http://localhost:8080/login")

	// ----------------- Static Files -----------------

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ----------------- Routes -----------------

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			db, err := database.InitDb()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			storedPassword, err := database.GetUser(db, username)
			if err != nil {
				http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
				return
			}

			if password != storedPassword {
				http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
				return
			}

			// set a cookie with the username
			cookie := http.Cookie{Name: "user", Value: username}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/forum", http.StatusSeeOther)
			fmt.Printf("Logged in user: %s\n", username)
			return
		}

		data := struct {
			User string
		}{
			User: "",
		}

		t := template.New("login")
		t = template.Must(t.ParseFiles("public/login.html", "public/template.html"))

		err := t.ExecuteTemplate(w, "login", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// http.ServeFile(w, r, "public/login.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			// add the user to the database
			db, err := database.InitDb()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			err = database.AddUser(db, username, password)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/forum", http.StatusSeeOther)
			fmt.Printf("A new user as register as: %s\n", username)
			return
		}

		data := struct {
			User string
		}{
			User: "",
		}

		t := template.New("register")
		t = template.Must(t.ParseFiles("public/register.html", "public/template.html"))

		err := t.ExecuteTemplate(w, "register", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//http.ServeFile(w, r, "public/register.html")
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// delete the cookie
		cookie := http.Cookie{Name: "user", Value: "", MaxAge: -1}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/forum", http.StatusSeeOther)
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// récupérer les valeurs du formulaire
			title := r.FormValue("title")
			content := r.FormValue("content")

			db, err := database.InitDb()
			if err != nil {
				http.Error(w, "Failed to initialize database", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			// ajouter la catégorie dans la base de données
			err = database.AddCategory(db, title, content)
			if err != nil {
				http.Error(w, "Failed to add the new category to the database", http.StatusInternalServerError)
				return
			}

			// rediriger vers la page d'accueil
			http.Redirect(w, r, "/forum", http.StatusSeeOther)
			return
		}

		// afficher le formulaire

		t := template.New("new_category")
		t = template.Must(t.ParseFiles("public/new_category.html", "public/template.html"))

		err := t.ExecuteTemplate(w, "new_category", struct{ User string }{User: ""})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// http.HandleFunc("/categories/new", func(w http.ResponseWriter, r *http.Request) {
	// 	err := templates.ExecuteTemplate(w, "new_category", struct{ User string }{User: ""})
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// })

	http.HandleFunc("/forum", func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("user")
		user := ""
		if err == nil {
			user = cookie.Value
		}

		data := struct {
			User string
		}{
			User: user,
		}

		tmpl := template.Must(template.ParseFiles("public/index.html"))
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Template forum is not executed", http.StatusInternalServerError)
		}

	})

	// ----------------- Listen And Serve -----------------
	http.ListenAndServe(":8080", nil)
}
