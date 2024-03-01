package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var (
	db          *sql.DB
	store       = sessions.NewCookieStore([]byte("Shin2")) // Changez la clé secrète
	sessionName = "session-name"
)

func main() {
	// Connexion à la base de données et gestion des routes
	initDB()
	defer db.Close()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/profil", profilHandler)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Utilisateur struct {
	Nom    string
	Prenom string
	Pseudo string
	Email  string
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/forum")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// L'utilisateur est connecté, affichez la page d'accueil
	tmpl := template.Must(template.ParseFiles("home.html")) // Assurez-vous d'avoir un fichier home.html avec votre contenu d'accueil
	tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Afficher le formulaire d'inscription s'il s'agit d'une requête GET
		tmpl := template.Must(template.ParseFiles("register.html"))
		tmpl.Execute(w, nil)
		return
	}

	// Récupérer les données du formulaire d'inscription
	nom := r.FormValue("nom")
	prenom := r.FormValue("prenom")
	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")
	motDePasse := r.FormValue("mot_de_passe")

	// Récupérer le fichier image téléversé
	file, handler, err := r.FormFile("image_profil")
	if err != nil {
		http.Error(w, "Erreur lors du téléversement de l'image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Enregistrer le fichier sur le serveur
	chemin := "/image" + handler.Filename
	fichier, err := os.Create(chemin)
	if err != nil {
		http.Error(w, "Erreur lors de la création du fichier", http.StatusInternalServerError)
		return
	}
	defer fichier.Close()
	io.Copy(fichier, file)

	// Insérer les données dans la base de données avec le chemin de l'image
	_, err = db.Exec("INSERT INTO utilisateurs (nom, prenom, pseudo, email, mot_de_passe, image_profil) VALUES (?, ?, ?, ?, ?, ?)", nom, prenom, pseudo, email, motDePasse, chemin)
	if err != nil {
		http.Error(w, "Erreur lors de l'inscription de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page de connexion après inscription réussie
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Afficher le formulaire de connexion s'il s'agit d'une requête GET
		tmpl := template.Must(template.ParseFiles("login.html")) // Assurez-vous d'avoir un fichier login.html avec le formulaire de connexion
		tmpl.Execute(w, nil)
		return
	}

	// Récupérer les données du formulaire de connexion
	pseudo := r.FormValue("pseudo")
	motDePasse := r.FormValue("mot_de_passe")

	// Récupérer le mot de passe correspondant au pseudo de la base de données
	var dbMotDePasse string
	err := db.QueryRow("SELECT mot_de_passe FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&dbMotDePasse)
	if err != nil {
		if err == sql.ErrNoRows {
			// Aucun utilisateur trouvé avec ce pseudo
			http.Error(w, "Nom d'utilisateur invalide", http.StatusUnauthorized)
		} else {
			// Erreur autre que sql.ErrNoRows
			http.Error(w, "Erreur lors de la connexion", http.StatusInternalServerError)
		}
		return
	}
	fmt.Println("Mot de passe récupéré de la base de données:", dbMotDePasse)

	// Comparer le mot de passe stocké avec celui fourni
	if motDePasse != dbMotDePasse {
		// Mot de passe invalide
		http.Error(w, "Mot de passe invalide", http.StatusUnauthorized)
		return
	}

	// Authentification réussie, définir la session comme authentifiée
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["pseudo"] = pseudo
	session.Save(r, w)

	// Rediriger l'utilisateur vers la page d'accueil après connexion réussie
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Réinitialiser la session
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = false
	session.Save(r, w)

	// Rediriger l'utilisateur vers la page de connexion
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations de l'utilisateur depuis la session
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}

	pseudo, ok := session.Values["pseudo"].(string)
	if !ok {
		http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
		return
	}

	// Récupérer les informations de l'utilisateur depuis la base de données
	var utilisateur Utilisateur
	err = db.QueryRow("SELECT  nom, prenom, email, pseudo  FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&utilisateur.Nom, &utilisateur.Prenom, &utilisateur.Pseudo, &utilisateur.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Afficher les informations de l'utilisateur sur la page HTML du profil
	tmpl := template.Must(template.ParseFiles("profil.html"))
	tmpl.Execute(w, utilisateur)
}
