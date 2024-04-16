package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
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
	http.HandleFunc("/topic", topicHandler)
	http.HandleFunc("/add-comment", addCommentHandler)
	http.HandleFunc("/thewitcher", thewitcher)
	http.HandleFunc("/ff7", finalFantasy7Handler)
	http.HandleFunc("/post", pageHandler)
	http.HandleFunc("/editopic", getMyTopicsHandler)
	http.HandleFunc("/delete-topic", deleteTopicHandler)
	http.HandleFunc("/edit-topic", editTopicHandler)
	http.HandleFunc("/zelda", zeldaHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/create-topic", createTopicHandler)
	http.HandleFunc("/propos", proposHandler)
	http.HandleFunc("/yakuza", yakuzaHandler)
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
	Image_profil string
}

type Sujet struct {
	ID        int
	Titre     string
	Contenu   string
	Auteur    string
	nomDuJeux string
}

type Message struct {
	ID       int
	Contenu  string
	Auteur   string
	CSSClass string
}

func init() {
	// Enregistrer le type map[string]bool pour éviter l'erreur de gob
	gob.Register(map[string]bool{})
	// Utiliser SecureCookie pour une meilleure sécurité
	store.Codecs = []securecookie.Codec{securecookie.New([]byte("your-secret-key"), nil)}
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
	pseudo, ok := session.Values["pseudo"].(string)
	if !ok {
		http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
		return
	}

	// Récupérer les informations de l'utilisateur depuis la base de données
	var utilisateur Utilisateur
	err := db.QueryRow("SELECT image_profil FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&utilisateur.Image_profil)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
		return
	}
	// Récupérer tous les sujets du forum depuis la base de données
	// Afficher la liste des sujets sur la page HTML du forum
	tmpl := template.Must(template.ParseFiles("forum.html"))
	tmpl.Execute(w, utilisateur)
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

	// Enregistrer l'image de profil sur le serveur
	file, handler, err := r.FormFile("image_profil")
	if err != nil {
		http.Error(w, "Erreur lors du téléchargement de l'image de profil", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Enregistrer le fichier sur le serveur avec un nom unique
	fileName := uuid.New().String() + filepath.Ext(handler.Filename)
	filePath := filepath.Join("Assets/profil", fileName)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erreur lors de l'enregistrement de l'image de profil sur le serveur", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copier le contenu du fichier dans le fichier de sortie
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Erreur lors de l'enregistrement de l'image de profil sur le serveur", http.StatusInternalServerError)
		return
	}

	// Insérer les données dans la base de données avec le nom de fichier de l'image de profil
	_, err = db.Exec("INSERT INTO utilisateurs (nom, prenom, pseudo, email, mot_de_passe, image_profil) VALUES (?, ?, ?, ?, ?, ?)", nom, prenom, pseudo, email, motDePasse, fileName)
	if err != nil {
		http.Error(w, "Erreur lors de l'inscription de l'utilisateur", http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["pseudo"] = pseudo
	session.Save(r, w)

	// Rediriger l'utilisateur vers la page d'accueil après l'inscription réussie
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
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pseudo, ok := session.Values["pseudo"].(string)
	if !ok {
		http.Error(w, "Utilisateur non connecté", http.StatusUnauthorized)
		return
	}

	// Récupérer les informations de l'utilisateur depuis la base de données
	var utilisateur Utilisateur
	err = db.QueryRow("SELECT nom, prenom, email, pseudo, image_profil FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&utilisateur.Nom, &utilisateur.Prenom, &utilisateur.Email, &utilisateur.Pseudo, &utilisateur.Image_profil)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Afficher les informations de l'utilisateur sur la page HTML du profil
	tmpl := template.Must(template.ParseFiles("profil.html"))
	tmpl.Execute(w, utilisateur)
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté à partir de la session
	session, _ := store.Get(r, sessionName)
	user, _ := session.Values["pseudo"].(string) // Supposons que le pseudo de l'utilisateur soit stocké dans la session

	// Récupérer l'identifiant du sujet à partir des paramètres de requête
	topicID := r.URL.Query().Get("id")

	// Récupérer les détails du sujet depuis la base de données
	var sujet Sujet
	err := db.QueryRow("SELECT titre, contenu, auteur, id, nomDuJeux FROM sujets WHERE id = ?", topicID).Scan(&sujet.Titre, &sujet.Contenu, &sujet.Auteur, &sujet.ID, &sujet.nomDuJeux)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des détails du sujet", http.StatusInternalServerError)
		return
	}

	// Récupérer tous les commentaires associés à ce sujet depuis la base de données
	commentaires, err := getCommentairesFromDB(topicID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des commentaires", http.StatusInternalServerError)
		return
	}

	// Parcourir les commentaires et déterminer le style CSS en fonction de l'auteur
	for i := range commentaires {
		if commentaires[i].Auteur == user {
			commentaires[i].CSSClass = "message-user"
		} else {
			commentaires[i].CSSClass = "message-other"
		}
	}

	// Déclaration des variables pour les sujets
	var topics []Sujet
	var templateFile string

	// Charger le template correspondant en fonction du nom du jeu et récupérer les sujets associés
	switch sujet.nomDuJeux {
	case "Final Fantasy":
		templateFile = "Jeux/Final-Fantasy/ff7_topics.html"
		topics, err = getSujetsByJeuFromDB("Final Fantasy")
	case "The witcher":
		templateFile = "Jeux/Thewitcher/Thewitcher3_topics.html"
		topics, err = getSujetsByJeuFromDB("The witcher")
	case "Zelda":
		templateFile = "Jeux/Zelda/Zelda BOTW_topics.html"
		topics, err = getSujetsByJeuFromDB("Zelda")
	case "Yakuza":
		templateFile = "Jeux/Yakuza/Yakuza_topics.html"
		topics, err = getSujetsByJeuFromDB("Yakuza")
	default:
		templateFile = "Jeux/Final-Fantasy/ff7_topics.html"
	}

	// Gérer les erreurs de récupération des sujets
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des sujets", http.StatusInternalServerError)
		return
	}

	// Afficher les détails du sujet et les commentaires associés sur la page HTML du sujet
	data := struct {
		Sujet        Sujet
		Commentaires []Message
		Topics       []Sujet
	}{
		Sujet:        sujet,
		Commentaires: commentaires,
		Topics:       topics,
	}

	// Charger le template HTML
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
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
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Récupérer les sujets associés à "The Witcher"
	sujets, err := getSujetsByJeuFromDB("The witcher")
	if err != nil {
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

func finalFantasy7Handler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Récupérer les sujets associés à "The Witcher"
	sujets, err := getSujetsByJeuFromDB("Final Fantasy")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des sujets", http.StatusInternalServerError)
		return
	}

	// Afficher les sujets sur la page HTML "The Witcher"
	tmpl := template.Must(template.ParseFiles("Jeux/Final-Fantasy/ff7.html"))
	err = tmpl.Execute(w, sujets)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du modèle HTML", http.StatusInternalServerError)
		return
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	// Charger la page HTML
	tmpl := template.Must(template.ParseFiles("post.html"))

	// Exécuter le modèle et écrire le résultat dans la réponse HTTP
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}
}

func getMyTopicsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Récupérer les sujets de l'utilisateur depuis la base de données
	sujets, err := getMyTopicsFromDB(pseudo)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de vos sujets", http.StatusInternalServerError)
		return
	}

	// Afficher les sujets de l'utilisateur sur la page HTML
	tmpl := template.Must(template.ParseFiles("ediTopic.html"))
	err = tmpl.Execute(w, sujets)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de vos sujets", http.StatusInternalServerError)
		return
	}
}

func getMyTopicsFromDB(pseudo string) ([]Sujet, error) {
	// Récupérer les sujets de l'utilisateur depuis la base de données en fonction de son pseudo
	rows, err := db.Query("SELECT id, titre FROM sujets WHERE auteur = ?", pseudo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créer une slice pour stocker les sujets de l'utilisateur
	var sujets []Sujet

	// Parcourir les résultats de la requête et ajouter les sujets à la slice
	for rows.Next() {
		var sujet Sujet
		err := rows.Scan(&sujet.ID, &sujet.Titre)
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

func deleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'identifiant du sujet à supprimer depuis les paramètres de requête
	topicID := r.URL.Query().Get("id")

	// Supprimer le sujet de la base de données en fonction de son identifiant
	_, err := db.Exec("DELETE FROM sujets WHERE id = ?", topicID)
	if err != nil {
		http.Error(w, "Erreur lors de la suppression du sujet", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page de ses sujets après la suppression
	http.Redirect(w, r, "/editopic", http.StatusSeeOther)
}

func editTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'identifiant du sujet à mettre à jour depuis les paramètres de requête
	topicID := r.URL.Query().Get("id")

	// Récupérer le nouveau titre du sujet depuis le formulaire de modification
	r.ParseForm()
	newTitle := r.FormValue("new_title")

	// Mettre à jour le titre du sujet dans la base de données
	_, err := db.Exec("UPDATE sujets SET titre = ? WHERE id = ?", newTitle, topicID)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour du titre du sujet", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page de ses sujets après la mise à jour
	http.Redirect(w, r, "/editopic", http.StatusSeeOther)
}

func zeldaHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sujets, err := getSujetsByJeuFromDB("Zelda")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des sujets", http.StatusInternalServerError)
		return
	}
	// Charger la page HTML
	tmpl := template.Must(template.ParseFiles("Jeux/Zelda/Zelda BOTW.html"))

	err = tmpl.Execute(w, sujets)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du modèle HTML", http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
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
	err = db.QueryRow("SELECT image_profil FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&utilisateur.Image_profil)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Afficher les informations de l'utilisateur sur la page HTML du profil
	tmpl := template.Must(template.ParseFiles("contact.html"))
	tmpl.Execute(w, utilisateur)
}

func proposHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
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
	err = db.QueryRow("SELECT image_profil FROM utilisateurs WHERE pseudo = ?", pseudo).Scan(&utilisateur.Image_profil)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Afficher les informations de l'utilisateur sur la page HTML du profil
	tmpl := template.Must(template.ParseFiles("a_propos.html"))
	tmpl.Execute(w, utilisateur)
}

func createTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si la méthode HTTP est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérifier si l'utilisateur est authentifié
	session, _ := store.Get(r, sessionName)
	if _, ok := session.Values["authenticated"].(bool); !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer les données du formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur lors de la lecture des données du formulaire", http.StatusInternalServerError)
		return
	}
	pseudo := session.Values["pseudo"].(string)
	// Récupérer le titre et le contenu du sujet à partir du formulaire
	titre := r.Form.Get("title")
	contenu := r.Form.Get("content")
	nomDuJeux := r.Form.Get("nomDuJeux")

	// Vérifier si les champs sont vides
	if titre == "" || contenu == "" {
		http.Error(w, "Veuillez remplir tous les champs du formulaire", http.StatusBadRequest)
		return
	}

	// Insérer le nouveau sujet dans la base de données
	_, err = db.Exec("INSERT INTO sujets (titre, contenu, nomDuJeux, auteur) VALUES (?, ?, ?, ?)", titre, contenu, nomDuJeux, pseudo)
	if err != nil {
		http.Error(w, "Erreur lors de la création du sujet", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page du forum après la création du sujet
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func yakuzaHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sujets, err := getSujetsByJeuFromDB("Yakuza")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des sujets", http.StatusInternalServerError)
		return
	}
	// Charger la page HTML
	tmpl := template.Must(template.ParseFiles("Jeux/Yakuza/Yakuza.html"))

	err = tmpl.Execute(w, sujets)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du modèle HTML", http.StatusInternalServerError)
		return
	}
}
