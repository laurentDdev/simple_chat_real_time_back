package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"instantmsg/models"
	"net/http"
	"strings"
	"time"
)

type userInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	Context *models.AppContext
}

func NewUserController(ctx *models.AppContext) *UserController {
	return &UserController{
		Context: ctx,
	}
}

type loginUser struct {
	Email    string
	Password string
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginu loginUser
	var u models.User

	err := json.NewDecoder(r.Body).Decode(&loginu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Vérification des informations d'identification dans la base de données
	err = uc.Context.DB.Where("email = ?", loginu.Email).First(&u).Error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Email not registered"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginu.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Password not match"))
		return
	}

	// Génération du token JWT
	token, err := uc.generateToken(&u)
	if err != nil {
		fmt.Printf("Erreur lors de la génération du token : %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := json.Marshal(struct {
		Email  string
		Pseudo string
		Id     uint
	}{
		Email:  u.Email,
		Pseudo: u.Pseudo,
		Id:     u.ID,
	})
	if err != nil {
		fmt.Printf("Erreur lors de la conversion au format json= %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Réponse réussie
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)

	w.Write(user)
}

func (uc *UserController) generateToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Hour)
	claims["id"] = user.ID
	key := []byte("powerdev") // Convertir la clé en []byte

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type user struct {
	Pseudo   string
	Mail     string
	Password string
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {

	var u user

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		r.Body.Close()
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Erreur lors du hashache du mot de passe %v\n", err)
		r.Body.Close()
		return
	}

	createdUser := models.User{
		Email:    u.Mail,
		Password: string(passwordHash),
		Pseudo:   u.Pseudo,
	}

	result := uc.Context.DB.Create(&createdUser)
	if result.Error != nil {
		fmt.Printf("Erreur lors de l'insertion de l'utilisateur %v\n", result.Error)
		if strings.HasPrefix(result.Error.Error(), "Error 1062 (23000): Duplicata du champ ") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Email already exist"))
			r.Body.Close()
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User inserted"))
	r.Body.Close()
	fmt.Printf("User inserted %v\n", result.RowsAffected)

}
