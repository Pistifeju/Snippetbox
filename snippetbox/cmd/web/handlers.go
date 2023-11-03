package main

import (
	"errors"
	"fmt"
	"github.com/Pistieju/snippetbox/internal/models"
	"github.com/Pistieju/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type snippetCreateForm struct {
	Title     string              `form:"title"`
	Content   string              `form:"content"`
	Expires   int                 `form:"expires"`
	Validator validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl.html", templateData)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl.html", templateData)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = snippetCreateForm{Expires: 365}
	app.render(w, http.StatusOK, "create.tmpl.html", templateData)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.Validator.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters")
	form.Validator.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.Validator.CheckField(validator.PermittedInt(form.Expires), "expires", "This field must equal 1, 7 or 365")

	if !form.Validator.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
