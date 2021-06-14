package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/taldrori/bookings/internal/config"
	"github.com/taldrori/bookings/internal/models"
)

var session *scs.SessionManager
var testapp config.Appconfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	//change to true in production
	testapp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testapp.InProduction

	testapp.Session = session

	app = &testapp
	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
