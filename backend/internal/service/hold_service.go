package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"gorm.io/gorm"
)

// HoldRepository 定义占座与库存扣减所需的数据访问能力。
type HoldRepository interface {
	ExistsActiveHoldTx(tx *gorm.DB, skuID, userID int64, now time.Time) (bool, error)
	CreateHoldTx(tx *gorm.DB, hold *domain.CabinHold) error
	AdjustInventoryTx(tx *gorm.DB, skuID int64, delta int, reason string) error
}

// CabinHoldService 负责舱位占座、过期时间计算与并发控制。
type CabinHoldService struct {
	repo    HoldRepository
	holdTTL time.Duration
	now     func() time.Time // 可注入，便于测试 TTL 过期场景
	locks   sync.Map
}

// NewCabinHoldService 创建占座服务，默认占座时长为 15 分钟。
func NewCabinHoldService(repo HoldRepository, holdTTL time.Duration) *CabinHoldService {
	if holdTTL <= 0 {
		holdTTL = 15 * time.Minute
	}
	return &CabinHoldService{repo: repo, holdTTL: holdTTL, now: time.Now}
}

// Hold 在无显式事务上下文时执行占座。
func (s *CabinHoldService) Hold(skuID int64, userID int64, qty int) bool {
	return s.HoldWithTx(nil, skuID, userID, qty)
}

// HoldWithTx 在指定事务中执行占座，并保证同用户同 SKU 的串行化处理。
func (s *CabinHoldService) HoldWithTx(tx *gorm.DB, skuID int64, userID int64, qty int) bool {
	if s.repo == nil || skuID <= 0 || userID <= 0 || qty <= 0 {
		return false
	}

	lockKey := fmt.Sprintf("%d:%d", skuID, userID)
	lock := s.loadLock(lockKey)
	lock.Lock()
	defer lock.Unlock()
	now := s.now()

	exists, err := s.repo.ExistsActiveHoldTx(tx, skuID, userID, now)
	if err != nil {
		return false
	}
	if exists {
		return true
	}

	if err := s.repo.AdjustInventoryTx(tx, skuID, -qty, "cabin_hold"); err != nil {
		return false
	}

	expiresAt := now.Add(s.holdTTL)
	if err := s.repo.CreateHoldTx(tx, &domain.CabinHold{
		CabinSKUID: skuID,
		UserID:     userID,
		Qty:        qty,
		ExpiresAt:  expiresAt,
	}); err != nil {
		_ = s.repo.AdjustInventoryTx(tx, skuID, qty, "cabin_hold_rollback")
		return false
	}

	return true
}

// loadLock 获取指定键对应的互斥锁，不存在时创建。
func (s *CabinHoldService) loadLock(key string) *sync.Mutex {
	v, _ := s.locks.LoadOrStore(key, &sync.Mutex{})
	return v.(*sync.Mutex)
}
