package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
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
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// test with none existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.ApartmentID = 50
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
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

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
