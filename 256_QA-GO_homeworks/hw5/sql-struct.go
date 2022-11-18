// Package hw5 is used for Homework 5
package hw5

import "time"

// DevicesTable is a struct for 'devices' table
type DevicesTable struct {
	ID        uint64    `db:"id"`
	Platform  string    `db:"platform"`
	UserId    uint64    `db:"user_id"`
	EnteredAt time.Time `db:"entered_at"`
	Removed   bool      `db:"removed"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DevicesEventsTable is a struct for 'devices_events' table
type DevicesEventsTable struct {
	ID        uint64    `db:"id"`
	DeviceID  uint64    `db:"device_id"`
	Type      uint8     `db:"type"`
	Status    uint8     `db:"status"`
	Payload   string    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// NotificationEventsTable is a struct for 'notification_events' table
type NotificationEventsTable struct {
	ID        uint64    `db:"id"`
	DeviceID  uint64    `db:"device_id"`
	Message   string    `db:"message"`
	Lang      uint8     `db:"lang"`
	Status    uint8     `db:"status"`
	Payload   string    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
