package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/doh-halle/ntarikoon-park/internal/config"
	"github.com/doh-halle/ntarikoon-park/internal/driver"
	"github.com/doh-halle/ntarikoon-park/internal/forms"
	"github.com/doh-halle/ntarikoon-park/internal/helpers"
	"github.com/doh-halle/ntarikoon-park/internal/models"
	"github.com/doh-halle/ntarikoon-park/internal/render"
	"github.com/doh-halle/ntarikoon-park/internal/repository"
	"github.com/doh-halle/ntarikoon-park/internal/repository/dbrepo"
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

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// send the data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})

}

//Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

//PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	err = errors.New("this is an error message")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName:   r.Form.Get("first_name"),
		LastName:    r.Form.Get("last_name"),
		Email:       r.Form.Get("email"),
		PhoneNumber: r.Form.Get("phone_number"),
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email", "phone_number")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

//Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

//Availability renders the search availability page and displays the form
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

//Post-Availability renders the search availability page and displays the form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJSON handles requests for availability and send JSON responds
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
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
	render.RenderTemplate(w, r, "anufor-embassy.page.tmpl", &models.TemplateData{})
}

//Mulang renders the mulang quarters page
func (m *Repository) Mulang(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "mulang-quarters.page.tmpl", &models.TemplateData{})
}

//Tawah renders the tawah house page
func (m *Repository) Tawah(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "tawah-house.page.tmpl", &models.TemplateData{})
}

//Anyere renders the anyere john suite page
func (m *Repository) Anyere(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "anyere-john-suite.page.tmpl", &models.TemplateData{})
}

//Ayafor renders the ayafor residence page
func (m *Repository) Ayafor(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "ayafor-residence.page.tmpl", &models.TemplateData{})
}

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
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
