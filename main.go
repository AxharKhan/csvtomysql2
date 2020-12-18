package main

import (
	"csvtomysql2/logic"
	"csvtomysql2/models"

	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl = template.Must(template.ParseGlob("form/*"))

//Index ..Default get method.
func Index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index", nil)
}

//Upload ..Upload data from CSV to Mysql.
func Upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	var errorstring string = ""
	file, header, er := r.FormFile("file") // where <<this>> is the controller and <<file>> the id of your form field
	if er != nil {
		errorstring = "Cannot Read the File please try again."
		tmpl.ExecuteTemplate(w, "error", errorstring)

		return
	}
	if file != nil {
		// get the filename
		fileName := header.Filename

		log.Println("file: " + fileName)

		s := strings.Split(fileName, ".")

		fileextension := s[1]

		if fileextension != "csv" {
			errorstring = "This File Type is not supported."
			tmpl.ExecuteTemplate(w, "error", errorstring)
		}
		// Parse the file
		r := csv.NewReader(file)

		persons := []models.Person{}

		for {

			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorstring = "Cannot Read the File please try again."
				tmpl.ExecuteTemplate(w, "error", errorstring)

			}

			var person models.Person
			var i int64

			person.FirstName = record[0]
			person.LastName = record[1]
			i, err = strconv.ParseInt(record[2], 10, 64)
			if err != nil {
				errorstring = "The data is not in the correct format. Please try another file. \n The correct Format is \"first_name,last_name,age,blood_group\" ."
				tmpl.ExecuteTemplate(w, "error", errorstring)
			}
			person.Age = i
			person.BloodGroup = record[3]

			persons = append(persons, person)
			logic.Insert(&person)

			if err != nil {
				errorstring = "An error occured while trying to store the data, Please try again."
				return
			}
			// csvdata += record[0] + "," + record[1] + "," + record[2] + "," + record[3] + ",\n"

		}

		var data models.Data
		data.Persons = persons
		data.Other = 1
		tmpl.ExecuteTemplate(w, "index", data)

	} else {
		errorstring = "Cannot Read the File please try again."
		tmpl.ExecuteTemplate(w, "error", errorstring)
	}
}

//ReadDB ..Reads Database Entries.
func ReadDB(w http.ResponseWriter, r *http.Request) {

	var persons []models.Person
	persons = logic.Read()
	var data models.Data
	data.Persons = persons
	data.Other = 0

	tmpl.ExecuteTemplate(w, "index", data)
}

func main() {
	for {
		err := logic.Create()
		if err != nil {
			log.Println("DB Connection failed. Retrying...")

		} else {
			break
		}
	}
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/readdb", ReadDB)
	http.ListenAndServe(":8080", nil)
}
