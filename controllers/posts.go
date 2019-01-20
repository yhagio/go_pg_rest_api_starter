package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_rest_pg_starter/middlewares"
	"go_rest_pg_starter/models"

	"github.com/gorilla/mux"
)

type Posts struct {
	ps models.PostService
}

type PostFormat struct {
	Title       string `schema:"title"`
	Description string `schema:"description"`
}

func NewPosts(ps models.PostService) *Posts {
	return &Posts{
		ps: ps,
	}
}

func (ps *Posts) Create(w http.ResponseWriter, r *http.Request) {
	user := middlewares.LookUpUserFromContext(r.Context())
	if user == nil {
		sendErrorResponse(w, http.StatusForbidden, "User not found.")
		return
	}

	// Get user input for creating a post
	var postFormat PostFormat
	err := json.NewDecoder(r.Body).Decode(&postFormat)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Cannot create a post.")
		return
	}

	post := models.Post{
		UserID:      user.ID,
		Title:       postFormat.Title,
		Description: postFormat.Description,
	}

	// Create a post
	err = ps.ps.Create(&post)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Failed to create a post.")
		return
	}

	setSuccessStatus(w, http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (p *Posts) GetOne(w http.ResponseWriter, r *http.Request) {
	post, err := p.getPostById(w, r)
	if err != nil {
		return
	}

	setSuccessStatus(w, http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (p *Posts) Update(w http.ResponseWriter, r *http.Request) {
	user := middlewares.LookUpUserFromContext(r.Context())
	if user == nil {
		sendErrorResponse(w, http.StatusForbidden, "User not found.")
		return
	}

	post, err := p.getPostById(w, r)
	if err != nil {
		return
	}

	if post.UserID != user.ID {
		sendErrorResponse(w, http.StatusForbidden, "You do not have permission.")
		return
	}

	var postFormat PostFormat
	err = json.NewDecoder(r.Body).Decode(&postFormat)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Cannot update the post.")
		return
	}

	updatingPost := models.Post{
		Title:       postFormat.Title,
		Description: postFormat.Description,
	}
	if updatingPost.Title != "" {
		post.Title = updatingPost.Title
	}
	if updatingPost.Description != "" {
		post.Description = updatingPost.Description
	}

	err = p.ps.Update(post)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Could not update the post.")
		return
	}

	setSuccessStatus(w, http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (p *Posts) Delete(w http.ResponseWriter, r *http.Request) {
	user := middlewares.LookUpUserFromContext(r.Context())
	if user == nil {
		sendErrorResponse(w, http.StatusForbidden, "User not found.")
		return
	}

	post, err := p.getPostById(w, r)
	if err != nil {
		return
	}

	if post.UserID != user.ID {
		sendErrorResponse(w, http.StatusForbidden, "You do not have permission.")
		return
	}
	err = p.ps.Delete(post.ID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Could not delete the post.")
		return
	}
}

// ------ Helper ------
func (p *Posts) getPostById(w http.ResponseWriter, r *http.Request) (*models.Post, error) {
	// Get :id from url id param, converted from string to int
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Invalid Post ID.")
		return nil, err
	}

	post, err := p.ps.GetOneById(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			sendErrorResponse(w, http.StatusNotFound, "Post not found.")
		default:
			sendErrorResponse(w, http.StatusInternalServerError, "Whoops! Something went wrong.")
		}
		return nil, err
	}

	return post, nil
}
