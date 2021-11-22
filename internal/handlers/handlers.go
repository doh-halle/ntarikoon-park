package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/doh-halle/ntarikoon-park/internal/config"
	"github.com/doh-halle/ntarikoon-park/internal/models"
	"github.com/doh-halle/ntarikoon-park/internal/render"
)

//Repo is the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

//About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//Perform some business logic
	stringMap := make(map[string]string)
	stringMap["Test"] = "Hello from Muea."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})

}

//Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
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
		log.Println(err)
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
