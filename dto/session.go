package dto

import (
	"time"
)

const DefaultSessionDurationMinutes int32 = 60 * 24

type Session struct {
	UserID           string
	ID               string
	Jwt              string
	Token            string
	StartedAt        time.Time
	LastAccessedAt   time.Time
	ExpiresAt        time.Time
	DeviceIPAddress  string
	DeviceUserAgent  string
	DeviceType       string
	IPAddressCity    string
	IPAddressCountry string
	User             User
}

type SessionListResponse struct {
	SessionID        string
	CurrentSession   bool
	LastAccessedAt   time.Time
	StartedAt        time.Time
	ExpiresAt        time.Time
	DeviceIPAddress  string
	DeviceUserAgent  string
	DeviceType       string
	IPAddressCity    string
	IPAddressCountry string
}

type SessionClaims struct {
	DeviceIPAddress  string
	DeviceUserAgent  string
	DeviceType       string
	IPAddressCity    string
	IPAddressCountry string
}
