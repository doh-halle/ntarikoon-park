package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/doh-halle/ntarikoon-park/internal/models"
)

//type postData struct {
//	key   string
//	value string
//}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"ae", "/anufor-embassy", "GET", http.StatusOK},
	{"mq", "/mulang-quarters", "GET", http.StatusOK},
	{"th", "/tawah-house", "GET", http.StatusOK},
	{"ajs", "/anyere-john-suite", "GET", http.StatusOK},
	{"ar", "/ayafor-residence", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"non-existent", "/other/types/of/reservations", "GET", http.StatusNotFound},
	// New routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	//{"post-sa", "/search-availability", "POST", []postData{
	//	{key: "start", value: "12-10-2021"},
	//	{key: "start", value: "12-10-2021"},
	//}, http.StatusOK},
	//{"post-sa-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "12-10-2021"},
	//	{key: "start", value: "12-10-2021"},
	//}, http.StatusOK},
	//{"post-mr", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value: "Justus"},
	//	{key: "last_name", value: "Fru"},
	//	{key: "email", value: "justfru@ayahoo.com"},
	//	{key: "phone_number", value: "+237-876-876-876"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		ApartmentID: 1,
		Apartment: models.Apartment{
			ID:            1,
			ApartmentName: "Anufor Embassy",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}
	// test with none existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.ApartmentID = 50
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

}

//func TestRepository_PostReservation(t *testing.T) {
//reqBody := "start_date=02-01-2035"
//reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=03-01-2035")
//reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Abraham")
//reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Neba")
//reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abran@pss.com")
//reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=+23755809358")
//reqBody = fmt.Sprintf("%s&%s", reqBody, "apartment_id=1")

//req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
//ctx := getCtx(req)
//req = req.WithContext(ctx)

//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

//rr := httptest.NewRecorder()

//handler := http.HandlerFunc(Repo.PostReservation)

//handler.ServeHTTP(rr, req)

//if rr.Code != http.StatusSeeOther {
//	t.Errorf("PostReservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
//}
//}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expeectedHTML      string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"bessem@orock.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"bih@mankon.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"su",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get thr url from test
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLocation.String())
			}
		}

		// checking for expected values in HTML
		if e.expeectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expeectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expeectedHTML)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
