package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles password hashing, JWT generation, and staff login.
type AuthService struct {
	staffRepo   *repository.StaffRepository
	jwtSecret   string
	expireHours int
}

// NewAuthService creates an AuthService.
func NewAuthService(staffRepo *repository.StaffRepository, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{staffRepo: staffRepo, jwtSecret: jwtSecret, expireHours: expireHours}
}

// HashPassword returns a bcrypt hash of the plaintext password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword returns true if hash matches password.
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateJWT signs a new JWT with HS256 for the given staffID and roles.
func GenerateJWT(staffID int64, roles []string, secret string, expireHours int) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", staffID),
		"roles": roles,
		"exp":   time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Login authenticates a staff member and returns a signed JWT and its expiry time.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, time.Time, error) {
	staff, err := s.staffRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}
	if staff.Status != 1 {
		return "", time.Time{}, errors.New("account disabled")
	}
	if !VerifyPassword(staff.PasswordHash, password) {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	// TODO: load roles from DB via staff_roles table (Sprint 2)
	roles := []string{"admin"}

	token, err := GenerateJWT(staff.ID, roles, s.jwtSecret, s.expireHours)
	if err != nil {
		return "", time.Time{}, err
	}
	expireAt := time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	return token, expireAt, nil
}

// GetProfile returns the staff domain object for the given sub claim string.
func (s *AuthService) GetProfile(ctx context.Context, subStr string) (*domain.Staff, error) {
	id, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid staff id")
	}
	return s.staffRepo.GetByID(ctx, id)
}
