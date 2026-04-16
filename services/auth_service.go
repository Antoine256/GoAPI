package services

import (
	"GoAPI/repository"
	"GoAPI/ressources"
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Login(dto ressources.LoginRequest, logger *zap.Logger) (ressources.TokenResponse, error) {
	// Récupère l'utilisateur par email
	user, err := repository.GetUserByName(dto.Name)
	if err != nil {
		logger.Error("Login - user not found", zap.String("name", dto.Name))
		return ressources.TokenResponse{}, errors.New("{'incorrect credential': 'name'}")
	}

	// Vérifie le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		logger.Error("Login - invalid password", zap.String("name", dto.Name))
		return ressources.TokenResponse{}, errors.New("{'incorrect credential': 'password'}")
	}

	// Génère les tokens
	tokens, err := GenerateTokensForUser(user)
	if err != nil {
		logger.Error("Login - failed to generate tokens", zap.String("name", dto.Name), zap.Error(err))
		return ressources.TokenResponse{}, err
	}

	return tokens, nil
}

func Register(dto ressources.RegisterRequest, logger *zap.Logger) (ressources.TokenResponse, error) {

	// Vérifie qu'il y a un Email (Si email, il faut faire une vérification sinon inutilisable !!!)
	// if strings.Trim(dto.Email, " ") == "" {

	// } else {
	// 	_, err := repository.GetUserByEmail(dto.Email)
	// 	if err == nil {
	// 		logger.Error("Register - email already exists", zap.String("email", dto.Email))
	// 		return ressources.TokenResponse{}, errors.New("{'incorrect credential': 'email'}")
	// 	}
	// }

	// Vérification de la validité du nom d'utilisateur
	_, err := repository.GetUserByName(dto.Name)
	if err == nil {
		logger.Error("Register - name already exists", zap.String("name", dto.Name))
		return ressources.TokenResponse{}, errors.New("{\"incorrect credential\": \"name\"}")
	}

	// Hash du mot de passe
	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Register - failed to hash password", zap.String("name", dto.Name), zap.Error(err))
		return ressources.TokenResponse{}, err
	}

	// Crée l'utilisateur
	user, err := repository.CreateUser(ressources.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: string(hashed),
		Role:     "user", // rôle par défaut
	})
	if err != nil {
		logger.Error("Register - failed to create user", zap.String("name", dto.Name), zap.Error(err))
		return ressources.TokenResponse{}, err
	}

	// Génère les tokens avec les infos du user
	tokens, err := GenerateTokensForUser(user)
	if err != nil {
		logger.Error("Register - failed to generate tokens", zap.String("name", dto.Name), zap.Error(err))
		return ressources.TokenResponse{}, err
	}

	return tokens, nil
}

func RefreshToken(refreshToken string, logger *zap.Logger) (ressources.TokenResponse, error) {
	// Vérifie et parse le refresh token
	claims, err := parseToken(refreshToken)
	if err != nil {
		logger.Error("RefreshToken - invalid token", zap.String("refresh_token", refreshToken))
		return ressources.TokenResponse{}, errors.New("token invalide")
	}

	userID := int(claims["user_id"].(float64))

	// Vérifie qu'il existe bien en base
	exists, err := repository.RefreshTokenExists(userID, refreshToken)
	if err != nil || !exists {
		logger.Error("RefreshToken - token not found in database", zap.Int("user_id", userID), zap.String("refresh_token", refreshToken))
		return ressources.TokenResponse{}, errors.New("token invalide ou expiré")
	}

	// Récupère les infos de l'utilisateur
	user, err := repository.GetUserByID(userID)
	if err != nil {
		logger.Error("RefreshToken - user not found", zap.Int("user_id", userID))
		return ressources.TokenResponse{}, errors.New("utilisateur non trouvé")
	}

	// Génère un nouvel access token
	accessToken, err := generateAccessToken(user)
	if err != nil {
		logger.Error("RefreshToken - failed to generate access token", zap.Int("user_id", userID), zap.Error(err))
		return ressources.TokenResponse{}, err
	}

	return ressources.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // on réutilise le même
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

func GenerateTokensForUser(user ressources.User) (ressources.TokenResponse, error) {
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return ressources.TokenResponse{}, err
	}
	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return ressources.TokenResponse{}, err
	}
	if err := repository.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return ressources.TokenResponse{}, err
	}

	return ressources.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}
