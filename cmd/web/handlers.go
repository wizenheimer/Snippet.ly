package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wizenheimer/snippet.ly/internal/models"
	"github.com/wizenheimer/snippet.ly/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title field can't be empty.")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content field can't be empty.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title field can't exceed 100 characters.")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "Expires field must be 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippet.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully.")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("user signup")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("user signup")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("user login")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("user login")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("user logout")
}
