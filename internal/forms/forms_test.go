package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form show valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/something", nil)

	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form show not valid when all required fields exist")
	}
}

func TestHas(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)
	if form.Has("a") {
		t.Error("Form show valid when checking for a missing field")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	r = httptest.NewRequest("POST", "/something", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	if !form.Has("a") {
		t.Error("Form show not valid when checking for a existing field")
	}
}

func TestMinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "123")
	r := httptest.NewRequest("POST", "/something", nil)
	r.PostForm = postedData
	form := New(postedData)
	if form.MinLength("a", 5) {
		t.Error("Form show that value of size 3 is good when required is 5")
	}
	if !form.MinLength("a", 2) {
		t.Error("Form show that value of size 3 is not good when required is 2")
	}
}

func TestIsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("NotEmail", "not_email")
	form := New(postedData)
	form.IsEmail("NotEmail")
	if form.Valid() {
		t.Error("Form show valid when checking an email which is not in the right format")
	}

	postedData = url.Values{}
	postedData.Add("Email", "a@a.com")
	form = New(postedData)
	form.IsEmail("Email")
	if !form.Valid() {
		t.Error("Form show not valid when checking an email when is in the right format")
	}
}
