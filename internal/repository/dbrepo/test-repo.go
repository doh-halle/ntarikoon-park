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

//AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

//AllNewReservations returns a slice of new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

//GetReservationByID returns one reservation by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var res models.Reservation

	return res, nil
}

//UpdateReservation updates a reservation in the database
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil
}

//DeletReservation deletes one reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {

	return nil
}

//UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

func (m *testDBRepo) AllApartments() ([]models.Apartment, error) {

	var apartments []models.Apartment

	return apartments, nil
}

//GetRestrictionsForApartmentByDate returns restrictions for an apartment by date range
func (m *testDBRepo) GetRestrictionsForApartmentByDate(apartmentID int, start, end time.Time) ([]models.ApartmentRestriction, error) {
	var restrictions []models.ApartmentRestriction
	return restrictions, nil

}

//InsertBlockforApartment inserts an apartment restriction
func (m *testDBRepo) InsertBlockforApartmentByID(id int, startDate time.Time) error {
	return nil
}

//DeleteBlockforApartment deletes an apartment restriction
func (m *testDBRepo) DeleteBlockforApartmentByID(id int) error {
	return nil
}
