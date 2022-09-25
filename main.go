package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	// memanggil package connection
	connection.DatabaseConnect()

	// route path folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// routing
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/form-project", formAddProject).Methods("GET")
	route.HandleFunc("/form-editproject/{index}", formEditProject).Methods("GET")
	route.HandleFunc("/detail-project/{index}", detailProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/edit-project", editProject).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")

	fmt.Println("server running on port 5050")
	// menjalankan server
	http.ListenAndServe("localhost:5050", route)
}

// membuat teplate (dto = data transformation object)
type Project struct {
	ID          int
	ProjectName string
	Description string
	StartDate   string
	EndDate     string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Description-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	// menampilkan data dari database

	data, _ := connection.Conn.Query(context.Background(), "SELECT id_project, project_name, description FROM tb_project ")
	// fmt.Println(data)

	var result []Project

	for data.Next() {
		var newdata = Project{}

		var err = data.Scan(&newdata.ID, &newdata.ProjectName, &newdata.Description)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, newdata)
	}

	resData := map[string]interface{}{
		"Projects": result,
	}
	fmt.Println(result)
	tmpl.Execute(w, resData)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Description-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/addproject.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

// mengambl data yang di push dari newproject dalam bentuk array
var dataProject = []Project{}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("ProjectName: " + r.PostForm.Get("ProjectName"))
	// fmt.Println("Description: " + r.PostForm.Get("description"))

	title := r.PostForm.Get("projectName")
	content := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startDate")

	// menyimpan data dari form ke object
	newProject := Project{
		ProjectName: title,
		Description: content,
		StartDate:   startDate,
	}

	// push data dari onject ke dataproject
	dataProject = append(dataProject, newProject)
	fmt.Println(dataProject)
	// redirect ke halaman index setelah button di klik
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Description-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/detailproject.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	// menangkap index params dari url
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// index, _ := mux.Vars(r)["index"]

	var projectDetail = Project{}

	// menampung data object dan menyesuaikan i dari data lopp dan index data param
	for i, data := range dataProject {
		if i == index {
			projectDetail = Project{
				ProjectName: data.ProjectName,
				Description: data.Description,
			}
		}
	}
	// menampilkan data
	data := map[string]interface{}{
		"Project": projectDetail,
	}

	tmpl.Execute(w, data)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Description-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Description-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/editproject.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	var projectEdit = Project{}

	// menampung data object dan menyesuaikan i dari data lopp dan index data param
	for i, data := range dataProject {
		if i == index {
			projectEdit = Project{
				ProjectName: data.ProjectName,
				Description: data.Description,
			}
		}
	}
	// menampilkan data
	data := map[string]interface{}{
		"Project": projectEdit,
	}
	fmt.Println(data)

	tmpl.Execute(w, data)
}
func editProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("projectName")
	content := r.PostForm.Get("description")
	startDate := r.PostForm.Get("startDate")

	// menyimpan data dari form ke object
	newProject := Project{
		ProjectName: title,
		Description: content,
		StartDate:   startDate,
	}

	// push data dari onject ke dataproject
	dataProject = append(dataProject, newProject)
	fmt.Println(dataProject)
	// redirect ke halaman index setelah button di klik
	http.Redirect(w, r, "/", http.StatusFound)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)
	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}
