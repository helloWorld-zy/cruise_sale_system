package domain

// BookingPassenger 表示预订与乘客之间的关联关系（多对多中间表）。
type BookingPassenger struct {
	ID          int64 `gorm:"primaryKey"` // 主键 ID
	BookingID   int64 `gorm:"index"`      // 关联的预订 ID
	PassengerID int64 `gorm:"index"`      // 关联的乘客 ID
}
