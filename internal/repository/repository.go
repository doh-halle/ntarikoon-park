package repository

import "github.com/doh-halle/ntarikoon-park/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertApartmentRestriction(r models.ApartmentRestriction) error
}
