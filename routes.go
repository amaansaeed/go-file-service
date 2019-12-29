package main

import (
	"encoding/json"
	"file-service/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
)

func (a *app) uploadFile(w http.ResponseWriter, r *http.Request) {
	incoming := "upload"
	dir := "temp"

	r.ParseMultipartForm(10 << 20)

	upload, handler, err := r.FormFile(incoming)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error retrieving file"))
		return
	}
	defer upload.Close()

	u := uuid.New()
	fn := u.String() + path.Ext(handler.Filename)

	f, err := os.OpenFile("./"+dir+"/"+fn, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error writing file"))
		return
	}
	defer f.Close()
	io.Copy(f, upload)

	id1, _ := uuid.NewRandom()
	id2, _ := uuid.NewRandom()

	var file models.File
	file.ID = id1.String()
	file.UserID = id2.String()
	file.Name = handler.Filename
	file.Filename = fn
	file.ContentType = handler.Header.Get("Content-Type")

	err = file.NewFile(a.DB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, _ := json.Marshal(file)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

// func (a *app) downloadFile(w http.ResponseWriter, r *http.Request) {

// }

func (a *app) getFilesByUser(w http.ResponseWriter, r *http.Request) {
	// type data struct {
	// 	UserID string `json:"userId"`
	// }
	var err error
	var f models.File
	// var d data

	err = json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if f.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Incomplete credentials"))
		return
	}

	files, err := f.GetFilesByUserId(a.DB, f.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, _ := json.Marshal(files)
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}
