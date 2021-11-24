package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"ae", "/anufor-embassy", "GET", []postData{}, http.StatusOK},
	{"mq", "/mulang-quarters", "GET", []postData{}, http.StatusOK},
	{"th", "/tawah-house", "GET", []postData{}, http.StatusOK},
	{"ajs", "/anyere-john-suite", "GET", []postData{}, http.StatusOK},
	{"ar", "/ayafor-residence", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-sa", "/search-availability", "POST", []postData{
		{key: "start", value: "12-10-2021"},
		{key: "start", value: "12-10-2021"},
	}, http.StatusOK},
	{"post-sa-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "12-10-2021"},
		{key: "start", value: "12-10-2021"},
	}, http.StatusOK},
	{"post-mr", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Justus"},
		{key: "last_name", value: "Fru"},
		{key: "email", value: "justfru@ayahoo.com"},
		{key: "phone_number", value: "+237-876-876-876"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
