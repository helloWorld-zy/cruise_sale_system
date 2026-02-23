package domain

type BookingPassenger struct {
	ID          int64 `gorm:"primaryKey"`
	BookingID   int64 `gorm:"index"`
	PassengerID int64 `gorm:"index"`
}
