package models

type Video struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SeriesID  uint   `gorm:"not null;index:idx_series_id" json:"series_id"`
	VideoName string `gorm:"size:255;not null;index:idx_video_name" json:"video_name"`
	VideoPath string `gorm:"size:255;not null" json:"video_path"`

	// Optional
	/*FileSize    int64     `gorm:"default:0" json:"file_size,omitempty"`
	Duration    int       `gorm:"default:0" json:"duration,omitempty"` // in seconds
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`*/

	// Virtual fields (not stored in DB)
	FileSizeFormatted string `gorm:"-" json:"file_size_formatted,omitempty"`
	DurationFormatted string `gorm:"-" json:"duration_formatted,omitempty"`
	FileExists        bool   `gorm:"-" json:"file_exists,omitempty"`
}
