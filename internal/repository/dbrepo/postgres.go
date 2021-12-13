package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/doh-halle/ntarikoon-park/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone_number, start_date, end_date, apartment_id, 
			created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.PhoneNumber,
		res.StartDate,
		res.EndDate,
		res.ApartmentID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

//InsertApartmentRestriction inserts an apartment restriction into the database
func (m *postgresDBRepo) InsertApartmentRestriction(r models.ApartmentRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into apartment_restrictions (start_date, end_date, apartment_id, reservation_id,
		created_at, updated_at, restriction_id)
		values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.ApartmentID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

//SearchAvailabilityByDatesByApartmentID returns true if availability exists for apartmentID and false if it doesnt
func (m *postgresDBRepo) SearchAvailabilityByDatesByApartmentID(start, end time.Time, apartmentID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select count(id) 
		from apartment_restrictions 
		where apartment_id = $1 and $2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, apartmentID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

//SearchAvailabilityForAllApartments returns a slice of available apartments if any, for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllApartments(start, end time.Time) ([]models.Apartment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var apartments []models.Apartment

	query := `
	select a.id, a.apartment_name
	from apartments a
	where a.id not in (select apartment_id from apartment_restrictions ar where $1 < ar.end_date and $2 > ar.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return apartments, err
	}

	for rows.Next() {
		var apartment models.Apartment
		err := rows.Scan(
			&apartment.ID,
			&apartment.ApartmentName,
		)
		if err != nil {
			return apartments, err
		}
		apartments = append(apartments, apartment)
	}

	if err = rows.Err(); err != nil {
		return apartments, err
	}

	return apartments, nil
}

// GetApartmentByID gets an apartment by ID
func (m *postgresDBRepo) GetApartmentByID(id int) (models.Apartment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var apartment models.Apartment

	query := `
		select id, apartment_name, created_at, updated_at from apartments where id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&apartment.ID,
		&apartment.ApartmentName,
		&apartment.CreatedAt,
		&apartment.UpdatedAt,
	)

	if err != nil {
		return apartment, err
	}

	return apartment, nil
}

//GetUserByID returns a user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
		from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}
	return u, nil
}

//UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

//Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

//AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone_number, r.start_date, 
		r.end_date, r.apartment_id, r.created_at, r.updated_at, r.processed,
		apm.id, apm.apartment_name
		from reservations r
		left join apartments apm on (r.apartment_id = apm.id)
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNumber,
			&i.StartDate,
			&i.EndDate,
			&i.ApartmentID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Apartment.ID,
			&i.Apartment.ApartmentName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

//AllNewReservations returns a slice of new reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone_number, r.start_date, 
		r.end_date, r.apartment_id, r.created_at, r.updated_at, 
		apm.id, apm.apartment_name
		from reservations r
		left join apartments apm on (r.apartment_id = apm.id)
		where processed = 0
		order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.PhoneNumber,
			&i.StartDate,
			&i.EndDate,
			&i.ApartmentID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Apartment.ID,
			&i.Apartment.ApartmentName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}
