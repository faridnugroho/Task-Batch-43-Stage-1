package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

type Project struct {
	Id                int
	Name              string
	Start_date        time.Time
	End_date          time.Time
	Format_start_date string
	Format_end_date   string
	Description       string
	Technologies      string
	Image             string
}

var Projects = []Project{
	{
		// Name:        "Dumbways Batch 43",
		// Start_date:  "2021-12-19",
		// End_date:    "2021-12-25",
		// Description: "Test",
	},
}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact-form", contactForm).Methods("GET")
	route.HandleFunc("/contact-form", addContactForm).Methods("POST")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/project/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/project", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	fmt.Println("Server is running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Format_start_date = each.Start_date.Format("2 January 2006")
		each.Format_end_date = each.End_date.Format("2 January 2006")

		result = append(result, each)
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func contactForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/contact-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addContactForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Name : " + r.PostForm.Get("name"))
	fmt.Println("Email : " + r.PostForm.Get("email"))
	fmt.Println("Phone : " + r.PostForm.Get("telp"))
	fmt.Println("Subject : " + r.PostForm.Get("subject"))
	fmt.Println("Message : " + r.PostForm.Get("message"))

	http.Redirect(w, r, "/contact-form", http.StatusMovedPermanently)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data": Data,
		"Id":   id,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	// start_date := r.PostForm.Get("start_date")
	// end_date := r.PostForm.Get("end_date")
	description := r.PostForm.Get("description")

	var newProject = Project{
		Name: name,
		// Start_date:  start_date,
		// End_date:    end_date,
		Description: description,
	}

	Projects = append(Projects, newProject)

	// fmt.Println(Projects)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	fmt.Println(id)

	Projects = append(Projects[:id], Projects[id+1:]...)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
