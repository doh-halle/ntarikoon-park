package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/doh-halle/ntarikoon-park/internal/config"
	"github.com/doh-halle/ntarikoon-park/internal/driver"
	"github.com/doh-halle/ntarikoon-park/internal/forms"
	"github.com/doh-halle/ntarikoon-park/internal/helpers"
	"github.com/doh-halle/ntarikoon-park/internal/models"
	"github.com/doh-halle/ntarikoon-park/internal/render"
	"github.com/doh-halle/ntarikoon-park/internal/repository"
	"github.com/doh-halle/ntarikoon-park/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

//Repo is the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})

}

//Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	apartment, err := m.DB.GetApartmentByID(res.ApartmentID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find an apartment!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Apartment.ApartmentName = apartment.ApartmentName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("02-01-2006")
	ed := res.EndDate.Format("02-01-2006")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cant get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.PhoneNumber = r.Form.Get("phone_number")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone_number")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "invalid form data", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.ApartmentRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		ApartmentID:   reservation.ApartmentID,
		ReservationID: newReservationID,
		RestrictionID: 2,
	}

	err = m.DB.InsertApartmentRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//send notifications - first to the guest
	htmlMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s, <br>
		This is to confirm your apartment reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "no-reply@ntarikoonpark.com",
		Subject:  "Reservation Confirmation - Ntarikoon Park",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//send notifications - to Property Owner
	htmlMessage = fmt.Sprintf(`
		<strong>Reservation Notification</strong><br>
		Dear Manager, <br>
		This is to confirm a reservation has been made for %s from %s to %s.
	`, reservation.Apartment.ApartmentName, reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"))

	msg = models.MailData{
		To:      "info@ntarikoonpark.com",
		From:    "info@ntarikoonpark.com",
		Subject: "Reservation Notification - Ntarikoon Park",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

//Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

//Availability renders the search availability page and displays the form
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

//Post-Availability renders the search availability page and displays the form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "02-01-2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	apartments, err := m.DB.SearchAvailabilityForAllApartments(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(apartments) == 0 {
		//No Availability

		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["apartments"] = apartments

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-apartment.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK          bool   `json:"ok"`
	Message     string `json:"message"`
	ApartmentID string `json:"apartment_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

//AvailabilityJSON handles requests for availability and send JSON responds
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "02-01-2006"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	apartmentID, _ := strconv.Atoi(r.Form.Get("apartment_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByApartmentID(startDate, endDate, apartmentID)
	resp := jsonResponse{
		OK:          available,
		Message:     "",
		StartDate:   sd,
		EndDate:     ed,
		ApartmentID: strconv.Itoa(apartmentID),
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Anufor renders the anufor embassy page
func (m *Repository) Anufor(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "anufor-embassy.page.tmpl", &models.TemplateData{})
}

//Mulang renders the mulang quarters page
func (m *Repository) Mulang(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "mulang-quarters.page.tmpl", &models.TemplateData{})
}

//Tawah renders the tawah house page
func (m *Repository) Tawah(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "tawah-house.page.tmpl", &models.TemplateData{})
}

//Anyere renders the anyere john suite page
func (m *Repository) Anyere(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "anyere-john-suite.page.tmpl", &models.TemplateData{})
}

//Ayafor renders the ayafor residence page
func (m *Repository) Ayafor(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "ayafor-residence.page.tmpl", &models.TemplateData{})
}

//ReservationSummary renders the reservation summary page - displaying reservation summary details
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get Item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("02-01-2006")
	ed := reservation.EndDate.Format("02-01-2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//ChooseApartment displays list of available apartments
func (m *Repository) ChooseApartment(w http.ResponseWriter, r *http.Request) {
	apartmentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.ApartmentID = apartmentID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

//ReserveApartment takes URL parameters, builts a sessional variable and takes user to make reservation page
func (m *Repository) ReserveApartment(w http.ResponseWriter, r *http.Request) {
	apartmentID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "02-01-2006"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	apartment, err := m.DB.GetApartmentByID(apartmentID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Apartment.ApartmentName = apartment.ApartmentName
	res.ApartmentID = apartmentID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return

	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in Successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})

}

//AdminNewReservation - returns all new reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

//AdminAllReservation - returns all reservations in admin tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{})

}
