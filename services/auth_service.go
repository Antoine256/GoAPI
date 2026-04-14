package services

import (
	"GoAPI/repository"
	"GoAPI/ressources"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Login(dto ressources.LoginDTO) (ressources.TokenResponseDTO, error) {
	// 1. Récupère l'utilisateur par email
	user, err := repository.GetUserByEmail(dto.Email)
	if err != nil {
		return ressources.TokenResponseDTO{}, errors.New("identifiants invalides")
	}

	// 2. Vérifie le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return ressources.TokenResponseDTO{}, errors.New("identifiants invalides")
	}

	// 3. Génère les tokens
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	// 4. Sauvegarde le refresh token en base
	if err := repository.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	return ressources.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func Register(dto ressources.RegisterDTO) (ressources.TokenResponseDTO, error) {
	// 1. Vérifie que l'email n'existe pas déjà
	_, err := repository.GetUserByEmail(dto.Email)
	if err == nil {
		return ressources.TokenResponseDTO{}, errors.New("email déjà utilisé")
	}

	// 2. Hash du mot de passe
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	// 3. Crée l'utilisateur
	user, err := repository.CreateUser(ressources.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: string(hashed),
		Role:     "user", // rôle par défaut
	})
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	// 4. Génère les tokens avec les infos du user
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}
	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	if err := repository.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	return ressources.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func RefreshToken(dto ressources.RefreshTokenDTO) (ressources.TokenResponseDTO, error) {
	// 1. Vérifie et parse le refresh token
	claims, err := parseToken(dto.RefreshToken)
	if err != nil {
		return ressources.TokenResponseDTO{}, errors.New("token invalide")
	}

	userID := int(claims["user_id"].(float64))

	// 2. Vérifie qu'il existe bien en base
	exists, err := repository.RefreshTokenExists(userID, dto.RefreshToken)
	if err != nil || !exists {
		return ressources.TokenResponseDTO{}, errors.New("token invalide ou expiré")
	}

	// 3. Récupère les infos de l'utilisateur
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return ressources.TokenResponseDTO{}, errors.New("utilisateur non trouvé")
	}

	// 4. Génère un nouvel access token
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return ressources.TokenResponseDTO{}, err
	}

	return ressources.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: dto.RefreshToken, // on réutilise le même
	}, nil
}

func Logout(refreshToken string) error {
	return repository.DeleteRefreshToken(refreshToken)
}

// --- Fonctions internes ---

func generateAccessToken(user ressources.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("méthode de signature invalide")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token invalide")
	}
	return token.Claims.(jwt.MapClaims), nil
}
