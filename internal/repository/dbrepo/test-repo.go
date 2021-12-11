package dbrepo

import (
	"errors"
	"time"

	"github.com/doh-halle/ntarikoon-park/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

//InsertApartmentRestriction inserts an apartment restriction into the database
func (m *testDBRepo) InsertApartmentRestriction(r models.ApartmentRestriction) error {
	return nil
}

//SearchAvailabilityByDatesByApartmentID returns true if availability exists for apartmentID and false if it doesnt
func (m *testDBRepo) SearchAvailabilityByDatesByApartmentID(start, end time.Time, apartmentID int) (bool, error) {
	return false, nil
}

//SearchAvailabilityForAllApartments returns a slice of available apartments if any, for a given date range
func (m *testDBRepo) SearchAvailabilityForAllApartments(start, end time.Time) ([]models.Apartment, error) {
	var apartments []models.Apartment
	return apartments, nil
}

// GetApartmentByID gets an apartment by ID
func (m *testDBRepo) GetApartmentByID(id int) (models.Apartment, error) {

	var apartment models.Apartment
	if id > 5 {
		return apartment, errors.New("something is not right - apartment needs to be cross-checked")
	}

	return apartment, nil
}

//GetUserByID gets users by ID
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User

	return u, nil
}

//UpdatesUser updates users
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

//Authenticate authenticates users using email and a password
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
