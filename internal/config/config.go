package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/doh-halle/ntarikoon-park/internal/models"
)

//AppConfig holds the application configurations settings
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
