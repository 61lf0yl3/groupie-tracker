package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Artist info
type Artist []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Relation info
type Relation struct {
	Index []Index `json:"index"`
}

// Index info
type Index struct {
	ID             int64               `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// About info
type About struct {
	ID           int
	Name         string
	Image        string
	Members      []string
	CreationDate int
	FirstAlbum   string
	RelationData Relation
}

// Search - what is returned in search
type Search struct {
	Artists   Artist
	RelationS Relation
}

var msgs = map[string]string{
	"album":    "Artsists / bands released their first album in ",
	"member":   " is a member of: ",
	"creation": "Artsists / bands created in ",
	"location": "Artsists / bands that have concerts in ",
	"date":     "Artsists / bands that have concerts in ",
}

var artists Artist
var relations Relation
var apiURL string = "https://groupietrackers.herokuapp.com/api/"
var err500 string = "500 Internal Server Error"
var err404 string = "404 This page not found"
var err400 string = "400 Bad Request"

func connectToAPI(s string) ([]byte, error) {
	resp, err := http.Get(apiURL + s)
	if err != nil {
		return nil, err
	}
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}
	return body, nil
}

func connnectParseAPI(w http.ResponseWriter, r *http.Request, api string) {
	body, err := connectToAPI(api)
	if err != nil {
		errorHandler(w, r, 500)
	}
	if api == "artists" {
		json.Unmarshal(body, &artists)
	} else if api == "relation" {
		json.Unmarshal(body, &relations)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if artists == nil {
		connnectParseAPI(w, r, "artists")
	}
	if r.URL.Path != "/" {
		errorHandler(w, r, 404)
		return
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, artists)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func redirectTo(id int, w http.ResponseWriter, r *http.Request) {
	page := "/about/" + strconv.Itoa(id)
	http.Redirect(w, r, page, http.StatusFound)
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	categ := r.FormValue("searchCateg")
	text := r.FormValue("searchText")
	if artists == nil {
		connnectParseAPI(w, r, "artists")
	}
	connnectParseAPI(w, r, "relation")

	creation, _ := strconv.Atoi(text)
	found := false
	var arr []int

	if categ == "album" || categ == "member" || categ == "artist" || categ == "creation" {
		for index, artist := range artists {
			if categ == "album" && artist.FirstAlbum == text {
				found = true
				arr = append(arr, index)
			} else if categ == "member" && contains(artist.Members, text) {
				found = true
				arr = append(arr, index)
			} else if categ == "artist" && artist.Name == text {
				found = true
				redirectTo(artist.ID, w, r)
			} else if categ == "creation" && artist.CreationDate == creation {
				found = true
				arr = append(arr, index)
			}
		}
	} else if categ == "location" {
		for index, relation := range relations.Index {
			if _, ok := relation.DatesLocations[text]; ok {
				found = true
				arr = append(arr, index)
			}
		}
	} else if categ == "date" {
		for index, relation := range relations.Index {
			for _, dates := range relation.DatesLocations {
				if contains(dates, text) {
					found = true
					arr = append(arr, index)
				}
			}
		}
	}

	t, _ := template.ParseFiles("templates/result.html")
	var msg string
	if categ == "member" {
		msg = text + msgs[categ]
	} else {
		msg = msgs[categ] + text + ": "
	}
	for i, index := range arr {
		if i == len(arr)-1 {
			msg += artists[index].Name
		} else {
			msg += artists[index].Name + ", "
		}
	}
	if !found {
		msg = "Nothing found for your request."
	}
	s := &About{Name: msg}
	t.Execute(w, s)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if artists == nil {
		connnectParseAPI(w, r, "artists")
	}
	id, _ := strconv.Atoi(r.URL.Path[len("/about/"):])
	if id >= 1 && id <= len(artists) {
		connnectParseAPI(w, r, "relation")
		id = id - 1
		artist := artists[id]
		about := &About{ID: id, Image: artist.Image, Name: artist.Name, Members: artist.Members, FirstAlbum: artist.FirstAlbum, CreationDate: artist.CreationDate, RelationData: relations}
		t, _ := template.ParseFiles("templates/about.html")
		t.Execute(w, about)
	} else {
		errorHandler(w, r, 400)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if artists == nil {
		connnectParseAPI(w, r, "artists")
	}
	searchArr := &Search{Artists: artists, RelationS: relations}
	t, _ := template.ParseFiles("templates/search.html")
	t.Execute(w, searchArr)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	t, _ := template.ParseFiles("templates/result.html")
	c := &About{}
	if status == http.StatusNotFound {
		c = &About{Name: err404}
	} else if status == 500 {
		c = &About{Name: err500}
	} else if status == 400 {
		c = &About{Name: err400}
	}
	t.Execute(w, c)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/about/", aboutHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/lookup/", lookupHandler)

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
