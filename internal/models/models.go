package models

import (
	"time"
)

//User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Apartment is the Apartment Model
type Apartment struct {
	ID            int
	ApartmentName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

//Restriction is the restrictions model
type Restriction struct {
	ID              int
	RestrictionName int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservation is the reservations model
type Reservation struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	StartDate   time.Time
	EndDate     time.Time
	ApartmentID int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Apartment   Apartment
}

//AppartmentRestrictions is the appartment_restriction model
type ApartmentRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	ApartmentID   int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Apartment     Apartment
	Reservation   Reservation
	Restriction   Restriction
}

//MailData holds email messages
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
