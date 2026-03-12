package domain

import "time"

type ContentTemplateKind string

const (
	ContentTemplateKindFeeNote       ContentTemplateKind = "fee_note"
	ContentTemplateKindBookingNotice ContentTemplateKind = "booking_notice"
)

type VoyageContentMode string

const (
	VoyageContentModeTemplate VoyageContentMode = "template"
	VoyageContentModeSnapshot VoyageContentMode = "snapshot"
)

type ContentTextItem struct {
	Text     string `json:"text"`
	Emphasis bool   `json:"emphasis,omitempty"`
}

type FeeNoteContent struct {
	Included []ContentTextItem `json:"included,omitempty"`
	Excluded []ContentTextItem `json:"excluded,omitempty"`
}

type BookingNoticeSection struct {
	Key   string            `json:"key"`
	Title string            `json:"title"`
	Items []ContentTextItem `json:"items,omitempty"`
}

type BookingNoticeContent struct {
	Sections []BookingNoticeSection `json:"sections,omitempty"`
}

type ContentTemplate struct {
	ID          int64               `gorm:"primaryKey" json:"id"`
	Name        string              `gorm:"size:120" json:"name"`
	Kind        ContentTemplateKind `gorm:"size:40;index" json:"kind"`
	Status      int16               `gorm:"default:1" json:"status"`
	ContentJSON string              `gorm:"column:content_json;type:text" json:"-"`
	Content     any                 `gorm:"-" json:"content,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
