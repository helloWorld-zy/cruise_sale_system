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

// AuthService 提供员工认证相关的业务逻辑，
// 包括密码哈希验证、JWT 令牌生成和员工登录功能。
type AuthService struct {
	staffRepo   *repository.StaffRepository // 员工数据仓储
	jwtSecret   string                      // JWT 签名密钥
	expireHours int                         // JWT 过期时间（小时）
}

// NewAuthService 创建认证服务实例，通过依赖注入传入员工仓储和 JWT 配置。
func NewAuthService(staffRepo *repository.StaffRepository, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{staffRepo: staffRepo, jwtSecret: jwtSecret, expireHours: expireHours}
}

// HashPassword 使用 bcrypt 算法对明文密码进行哈希处理。
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword 验证密码哈希与明文密码是否匹配。
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateJWT 使用 HS256 算法签发 JWT 令牌，包含员工 ID 和角色信息。
func GenerateJWT(staffID int64, roles []string, secret string, expireHours int) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", staffID),                                    // 主题：员工 ID
		"roles": roles,                                                         // 角色列表
		"exp":   time.Now().Add(time.Duration(expireHours) * time.Hour).Unix(), // 过期时间
		"iat":   time.Now().Unix(),                                             // 签发时间
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Login 验证员工凭据并返回签名后的 JWT 令牌及其过期时间。
// 验证流程：查找用户 → 检查账户状态 → 验证密码 → 签发令牌。
func (s *AuthService) Login(ctx context.Context, username, password string) (string, time.Time, error) {
	// 根据用户名查找员工
	staff, err := s.staffRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}
	// 检查账户是否已启用
	if staff.Status != 1 {
		return "", time.Time{}, errors.New("account disabled")
	}
	// 验证密码
	if !VerifyPassword(staff.PasswordHash, password) {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	// TODO: 从数据库 staff_roles 表加载角色（Sprint 2）
	roles := []string{"admin"}

	// 签发 JWT 令牌
	token, _ := GenerateJWT(staff.ID, roles, s.jwtSecret, s.expireHours)
	expireAt := time.Now().Add(time.Duration(s.expireHours) * time.Hour)
	return token, expireAt, nil
}

// GetProfile 根据 JWT 中的 sub 声明（员工 ID 字符串）查询员工信息。
func (s *AuthService) GetProfile(ctx context.Context, subStr string) (*domain.Staff, error) {
	id, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid staff id")
	}
	return s.staffRepo.GetByID(ctx, id)
}
