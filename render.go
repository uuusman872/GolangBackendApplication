package main

import (
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, filename string, data *templateData) {
	if app.tp == nil {
		http.Error(w, "tamplate rendering engine is not set", http.StatusInternalServerError)
		return
	}
	mergedData := app.defaultTemplateData(data, r)
	app.tp.Render(w, filename, mergedData)
}

func (app *application) defaultTemplateData(data *templateData, r *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}
	data.Flash = app.session.PopString(r, "flash")
	data.IsAuthenticated = app.IsAuthenticated(r)
	return data
}
