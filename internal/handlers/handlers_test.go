package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taldrori/bookings/internal/models"
)

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "Get", http.StatusOK},
	{"about", "/about", "Get", http.StatusOK},
	{"jq", "/jonins-quarters", "Get", http.StatusOK},
	{"hs", "/hokages-suite", "Get", http.StatusOK},
	{"search-availability", "/search-availability", "Get", http.StatusOK},
	{"contact", "/contact", "Get", http.StatusOK},
	// {"post-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "01/01/2020"},
	// 	{key: "end", value: "01/02/2020"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "01/01/2020"},
	// 	{key: "end", value: "01/02/2020"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Tal"},
	// 	{key: "last_name", value: "Drori"},
	// 	{key: "email", value: "t@t.com"},
	// 	{key: "phone", value: "123456789"},
	// }, http.StatusOK},
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
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Jonin's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where error happened when retrieving room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	reservation.RoomID = 100

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Jonin's Quarters",
		},
	}

	reqBody := "start_date=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Tal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// no reservation in context
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code when no reservation on context: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// missing test post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code when no form body: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid start_date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Tal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code when start_date is invalid: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end_date
	reqBody = "start_date=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Tal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code when end_date is invalid: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid room_id
	reqBody = "start_date=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Tal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code when room_id is invalid: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test if form isn't valid (first name is less then 3 letters)
	reqBody = "start_date=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=T")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned wrong response code invalid data: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for faliuare of insert to reservation
	reqBody = "start_date=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Tal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler failed when trying to insert to reservation: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for faliuare of insert to restriction
	reqBody = "start_date=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=02/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Tal")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Drori")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tal@drori.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555555555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler failed when trying to insert to restriction: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	// create our request body
	reqBody := "start=01/01/2050"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=01/02/2050")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.AvailabilityJson)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	// this time we want to parse JSON and get the expected response
	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date > 2049-12-31, we expect no availability
	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	/*****************************************
	// second case -- rooms not available
	*****************************************/
	// create our request body
	reqBody = "start=01/01/2040"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=01/02/2040")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJson)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if !j.OK {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}

	/*****************************************
	// third case -- no request body
	*****************************************/
	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	// get the context with session
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJson)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Internal server error" {
		t.Error("Got availability when request body was empty")
	}

	/*****************************************
	// fourth case -- database error
	*****************************************/
	// create our request body
	reqBody = "start=01/01/2060"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=01/02/2060")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailabilityJson)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// this time we want to parse JSON and get the expected response
	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Error querying database" {
		t.Error("Got availability when simulating database error")
	}
}

func getCTX(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
