package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/taldrori/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.Appconfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// as wanted
	default:
		t.Error(fmt.Sprintf("type is not chi.Mux, but is %T", v))
	}
}
