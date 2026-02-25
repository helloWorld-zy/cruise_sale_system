package domain

// BookingPassenger 表示预订与乘客之间的关联关系。
type BookingPassenger struct {
	ID          int64 `gorm:"primaryKey"`
	BookingID   int64 `gorm:"index"`
	PassengerID int64 `gorm:"index"`
}
