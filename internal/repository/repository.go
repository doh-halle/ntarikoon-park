package repository

import (
	"time"

	"github.com/doh-halle/ntarikoon-park/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertApartmentRestriction(r models.ApartmentRestriction) error
	SearchAvailabilityByDatesByApartmentID(start, end time.Time, apartmentID int) (bool, error)
	SearchAvailabilityForAllApartments(start, end time.Time) ([]models.Apartment, error)
	GetApartmentByID(id int) (models.Apartment, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
}
