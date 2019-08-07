package util

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"html/template"
	"io"
)

type TemplateRegistry struct {
	Templates map[string]*template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	invalidSessionError := CheckForValidSession(c)

	dateSession, _ := session.Get("date-session", c)

	if invalidSessionError == nil {
		if data != nil {
			data.(map[string]interface{})["sessionData"] = dateSession
			data.(map[string]interface{})["warningFlashMessages"] = dateSession.Flashes("warningFlashMessages")
			data.(map[string]interface{})["errorFlashMessages"] = dateSession.Flashes("errorFlashMessages")
			data.(map[string]interface{})["infoFlashMessages"] = dateSession.Flashes("infoFlashMessages")
			data.(map[string]interface{})["successFlashMessages"] = dateSession.Flashes("successFlashMessages")

		}
	}

	tmpl, ok := t.Templates[name]

	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}

	saveFailed := dateSession.Save(c.Request(), c.Response())

	if saveFailed != nil {
		return saveFailed
	}

	return tmpl.ExecuteTemplate(w, "base.html", data)
}

