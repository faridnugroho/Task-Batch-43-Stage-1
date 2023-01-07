package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = MetaData{
	Title: "Personal Web",
}

type Project struct {
	Id                int
	Title             string
	Start_date        time.Time
	End_date          time.Time
	Format_start_date string
	Format_end_date   string
	Duration          string
	Description       string
	Technologies      []string
	Images            string
	Author            string
	IsLogin           bool
}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

var Projects = []Project{}

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnection()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	route.HandleFunc("/home", home).Methods("GET").Name("home")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/project", middleware.UploadFile(addProject)).Methods("POST")
	route.HandleFunc("/project/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/contact-form", contactForm).Methods("GET")

	route.HandleFunc("/update-project/{id}", updateProject).Methods("GET")
	route.HandleFunc("/update-project/", formUpdate).Methods("POST")

	route.HandleFunc("/register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}
	Data.FlashData = strings.Join(flashes, "")

	rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, title, start_date, end_date, description, technologies, images, tb_user.name as author FROM tb_project LEFT JOIN tb_user ON tb_user.id = tb_project.author_id ORDER BY id DESC")

	var result []Project
	for rows.Next() {
		var each = Project{}

		var err = rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Images, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Duration = getDistanceTime(each.Start_date, each.End_date)

		if session.Values["IsLogin"] != true {
			each.IsLogin = false
		} else {
			each.IsLogin = session.Values["IsLogin"].(bool)
		}

		result = append(result, each)
	}

	respData := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	start_date := r.PostForm.Get("start_date")
	end_date := r.PostForm.Get("end_date")
	technologies := r.Form["tech"]

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	author := session.Values["Id"].(int)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_project(title, description, start_date, end_date, technologies, images, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", title, description, start_date, end_date, technologies, image, author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/project-update.html")

	result := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT title, start_date,end_date, description, technologies, images FROM tb_project WHERE id=$1", id).Scan(&result.Title, &result.Start_date, &result.End_date, &result.Description, &result.Technologies, &result.Images)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	result.Id = id

	resp := map[string]interface{}{
		"Data":   Data,
		"Result": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func formUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	id := r.PostForm.Get("id")
	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	start_date := r.PostForm.Get("start_date")
	end_date := r.PostForm.Get("end_date")
	technologies := r.Form["tech"]

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	author := session.Values["Id"].(int)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title=$1, description=$2, start_date=$3, end_date=$4, technologies=$5, images=$6, author_id=$7 WHERE id=$8", title, description, start_date, end_date, technologies, image, author, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
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

	ProjectDetail := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date,end_date, description, technologies, images FROM tb_project WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Images)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ProjectDetail.Duration = getDistanceTime(ProjectDetail.Start_date, ProjectDetail.End_date)

	ProjectDetail.Format_start_date = ProjectDetail.Start_date.Format("2 January 2006")
	ProjectDetail.Format_end_date = ProjectDetail.End_date.Format("2 January 2006")

	resp := map[string]interface{}{
		"Data":    Data,
		"Project": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func contactForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("views/contact-form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func getDistanceTime(Start_date, End_date time.Time) string {
	duration := End_date.Sub(Start_date)

	year := int(duration.Hours() / (12 * 30 * 24))
	month := int(duration.Hours() / (30 * 24))
	week := int(duration.Hours() / (7 * 24))
	day := int(duration.Hours() / 24)

	if year != 0 {
		return strconv.Itoa(year) + " year ago"
	}
	if month != 0 {
		return strconv.Itoa(month) + " month ago"
	}
	if week != 0 {
		return strconv.Itoa(week) + " week ago"
	}
	if day != 0 {
		return strconv.Itoa(day) + " days ago"
	}
	return "today"
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.AddFlash("Successfully register!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// cookie = storing data
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800 // 1 jam = 3600 detik | 3 jam = 10800

	session.AddFlash("successfully login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout.")
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1 // gak boleh kurang dari 0
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
