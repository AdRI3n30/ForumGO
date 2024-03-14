package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

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

	http.HandleFunc("/", forumHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/profil", profilHandler)
	http.HandleFunc("/create-topic", createTopicHandler)
	http.HandleFunc("/topic", topicHandler)
	http.HandleFunc("/add-comment", addCommentHandler)
	http.HandleFunc("/thewitcher", thewitcher)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/Assets/", http.StripPrefix("/Assets/", http.FileServer(http.Dir("Assets"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Utilisateur struct {
	Nom          string
	Prenom       string
	Pseudo       string
	Email        string
	Image_profil sql.NullString
}

type Sujet struct {
	ID        int
	Titre     string
	Contenu   string
	Auteur    string
	nomDuJeux string
}

type Message struct {
	ID      int
	Contenu string
	Auteur  string
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

func forumHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Récupérer tous les sujets du forum depuis la base de données
	sujets, err := getSujetsFromDB()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des sujets du forum", http.StatusInternalServerError)
		return
	}

	// Afficher la liste des sujets sur la page HTML du forum
	tmpl := template.Must(template.ParseFiles("forum.html"))
	tmpl.Execute(w, sujets)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Afficher le formulaire d'inscription s'il s'agit d'une requête GET
		tmpl := template.Must(template.ParseFiles("register.html")) // Assurez-vous d'avoir un fichier register.html avec le formulaire d'inscription
		tmpl.Execute(w, nil)
		return
	}

	// Récupérer les données du formulaire d'inscription
	nom := r.FormValue("nom")
	prenom := r.FormValue("prenom")
	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")
	motDePasse := r.FormValue("mot_de_passe")

	// Insérer les données dans la base de données avec le mot de passe en texte brut
	_, err := db.Exec("INSERT INTO utilisateurs (nom, prenom, pseudo, email, mot_de_passe) VALUES (?, ?, ?, ?, ?)", nom, prenom, pseudo, email, motDePasse)
	if err != nil {
		http.Error(w, "Erreur lors de l'inscription de l'utilisateur", http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["pseudo"] = pseudo
	session.Save(r, w)
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
	err = db.QueryRow("SELECT nom, prenom, email, pseudo, image_profil FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&utilisateur.Nom, &utilisateur.Prenom, &utilisateur.Pseudo, &utilisateur.Email, &utilisateur.Image_profil)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Afficher les informations de l'utilisateur sur la page HTML du profil
	tmpl := template.Must(template.ParseFiles("profil.html"))
	tmpl.Execute(w, utilisateur)
}

func getSujetsFromDB() ([]Sujet, error) {
	// Récupérer tous les sujets du forum depuis la base de données
	rows, err := db.Query("SELECT id, titre, auteur FROM sujets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créer une slice pour stocker les sujets
	var sujets []Sujet

	// Parcourir les résultats de la requête et ajouter les sujets à la slice
	for rows.Next() {
		var sujet Sujet
		err := rows.Scan(&sujet.ID, &sujet.Titre, &sujet.Auteur)
		if err != nil {
			return nil, err
		}
		sujets = append(sujets, sujet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sujets, nil
}

func createTopicHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	titre := r.Form.Get("title")
	contenu := r.Form.Get("content")
	nomDuJeux := r.FormValue("nomDuJeux")

	// Récupérer le pseudo de l'utilisateur depuis la session
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

	// Insérer le nouveau sujet dans la base de données avec le pseudo de l'utilisateur comme auteur
	_, err = db.Exec("INSERT INTO sujets (titre, contenu, auteur,nomDuJeux) VALUES (?, ?, ?, ?)", titre, contenu, pseudo, nomDuJeux)
	if err != nil {
		http.Error(w, "Error creating topic", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page d'accueil après création du sujet
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'identifiant du sujet à partir des paramètres de requête
	topicID := r.URL.Query().Get("id")
	fmt.Println(topicID)

	// Récupérer les détails du sujet depuis la base de données
	var sujet Sujet
	err := db.QueryRow("SELECT titre, contenu, auteur, id FROM sujets WHERE id = ?", topicID).Scan(&sujet.Titre, &sujet.Contenu, &sujet.Auteur, &sujet.ID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des détails du sujet", http.StatusInternalServerError)
		return
	}

	// Log des détails du sujet récupérés depuis la base de données
	log.Printf("Détails du sujet: %v", sujet)

	// Récupérer tous les commentaires associés à ce sujet depuis la base de données
	commentaires, err := getCommentairesFromDB(topicID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des commentaires", http.StatusInternalServerError)
		return
	}

	// Afficher les détails du sujet et les commentaires associés sur la page HTML du sujet
	tmpl := template.Must(template.ParseFiles("topic.html"))
	data := struct {
		Sujet        Sujet
		Commentaires []Message
	}{
		Sujet:        sujet,
		Commentaires: commentaires,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page du sujet", http.StatusInternalServerError)
		return
	}
}
func getCommentairesFromDB(topicID string) ([]Message, error) {
	// Récupérer tous les commentaires associés à ce sujet depuis la base de données
	rows, err := db.Query("SELECT id, contenu, auteur FROM messages WHERE sujet_id = ?", topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créer une slice pour stocker les commentaires
	var commentaires []Message

	// Parcourir les résultats de la requête et ajouter les commentaires à la slice
	for rows.Next() {
		var commentaire Message
		err := rows.Scan(&commentaire.ID, &commentaire.Contenu, &commentaire.Auteur)
		if err != nil {
			return nil, err
		}
		commentaires = append(commentaires, commentaire)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentaires, nil
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les données du formulaire de commentaire
	contenu := r.FormValue("contenu")
	sujetID := r.FormValue("sujet_id") // Récupérer l'ID du sujet depuis le formulaire

	// Récupérer le pseudo de l'utilisateur depuis la session
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

	// Insérer le nouveau commentaire dans la base de données avec le pseudo de l'utilisateur comme auteur
	_, err = db.Exec("INSERT INTO messages (contenu, auteur, sujet_id) VALUES (?, ?, ?)", contenu, pseudo, sujetID)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page du sujet après l'ajout du commentaire
	http.Redirect(w, r, "/topic?id="+sujetID, http.StatusSeeOther)
}

func thewitcher(w http.ResponseWriter, r *http.Request) {
	// Récupérer les sujets associés à "The Witcher"
	sujets, err := getSujetsByJeuFromDB("The witcher")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erreur lors de la récupération des sujets", http.StatusInternalServerError)
		return
	}

	// Afficher les sujets sur la page HTML "The Witcher"
	tmpl := template.Must(template.ParseFiles("Jeux/Thewitcher/TheWitcher.html"))
	err = tmpl.Execute(w, sujets)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du modèle HTML", http.StatusInternalServerError)
		return
	}
}

func getSujetsByJeuFromDB(jeu string) ([]Sujet, error) {
	// Exécuter la requête SQL pour récupérer les sujets en fonction du nom du jeu
	rows, err := db.Query("SELECT id, titre, contenu, auteur, nomDuJeux FROM sujets WHERE nomDuJeux = ?", jeu)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les lignes de résultats et construire une liste de sujets
	var sujets []Sujet
	for rows.Next() {
		var sujet Sujet
		err := rows.Scan(&sujet.ID, &sujet.Titre, &sujet.Contenu, &sujet.Auteur, &sujet.nomDuJeux)
		if err != nil {
			return nil, err
		}
		sujets = append(sujets, sujet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sujets, nil
}
