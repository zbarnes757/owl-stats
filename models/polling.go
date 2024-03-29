package models

import "time"

// Polling is used to keep track of polling tasks that need to be done
type Polling struct {
	Base
	LastRan time.Time
	NextRun time.Time
	Title   string `gorm:"unique;not null"`
}

// Update a polling with the given data
func (p *Polling) Update(data Polling) {
	db.Model(p).Updates(data)
}

// GetPolling finds a polling by title
func GetPolling(title string) *Polling {
	polling := &Polling{}

	if db.Table("pollings").Where("title = ?", title).First(polling).RecordNotFound() {
		return nil
	}

	return polling
}
