package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/iDos27/order-management/auth-service/internal/database"
	"github.com/iDos27/order-management/auth-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	db *database.DB
}

func NewAuthMiddleware(db *database.DB) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

func (m *AuthMiddleware) ValidateTokenEndpoint(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Authorization header is missing",
		})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Invalid Authorization header format",
		})
		return
	}

	tokenString := tokenParts[1]
	claims, err := m.validateJWTToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Invalid token: " + err.Error(),
		})
		return
	}

	var user models.User
	err = m.db.QueryRow(`
		SELECT id, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, claims["user_id"]).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"claims":  claims,
	})
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := tokenParts[1]
		claims, err := m.validateJWTToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Zapisanie w sesji
		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		c.Set("claims", claims)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exist := c.Get("role")
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Role not found",
				"message": "Role info not found",
			})
			return
		}

		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":          "insufficient_permissions",
			"message":        "User role '" + userRole + "' is not authorized for this resource",
			"required_roles": allowedRoles,
			"user_role":      userRole,
		})
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole(models.RoleAdmin)
}

func (m *AuthMiddleware) RequireAdminOrEmployee() gin.HandlerFunc {
	return m.RequireRole(models.RoleAdmin, models.RoleEmployee)
}

func (m *AuthMiddleware) validateJWTToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret-key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func GetCurrentUser(c *gin.Context) (*models.CurrentUser, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, jwt.ErrTokenNotValidYet
	}
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	userIDFloat, ok := userID.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user type: expected float64, got %T", userID)
	}

	return &models.CurrentUser{
		ID:    int(userIDFloat),
		Email: email.(string),
		Role:  role.(string),
	}, nil
}
