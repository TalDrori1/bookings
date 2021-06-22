package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/taldrori/bookings/internal/models"
)

// AppConfig holds the application config
type Appconfig struct {
	UseChache     bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
