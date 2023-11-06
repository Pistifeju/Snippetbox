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

type userSignUpForm struct {
	Name      string              `form:"name"`
	Email     string              `form:"email"`
	Password  string              `form:"password"`
	Validator validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email     string              `form:"email"`
	Password  string              `form:"password"`
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
	form.Validator.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = userSignUpForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", templateData)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.Validator.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.Validator.CheckField(validator.MatchesPattern(form.Email, validator.EmailRX), "email", "This field is not a valid email address")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.Validator.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Validator.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Validator.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", templateData)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.Validator.CheckField(validator.MatchesPattern(form.Email, validator.EmailRX), "email", "This field is not a valid email address")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Validator.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Validator.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "You've been logged in successfully!")
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
