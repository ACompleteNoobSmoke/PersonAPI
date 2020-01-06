package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*Payload has the a map containing everyone
added to the map and how it displays when requested*/
type Payload struct {
	Everyone People
}

//Person represents a person structure
type Person struct {
	Name     string `json:"name"`
	Age      int64
	FavColor string `json:"favorite_color"`
}

//People is a map containing each person that will be created
type People map[string]*Person

//p is a global map where newPeople will be stored
var p = make(map[string]*Person)

//pay is global variable to keep updating when more people are created for the payload
var pay Payload

//peopleRequest represents the function that handles the request user made to "/people"
func peopleRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getEveryone(w, r)
		break
	case http.MethodPost:
		inputPerson(w, r)
		break
	default:
		fmt.Fprintf(w, "Incorrect Method")
		break
	}
}

//inputPerson uses POST method to create a new person
func inputPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsondecode := json.NewDecoder(r.Body)

	peeps := &Person{}
	err := jsondecode.Decode(peeps)
	if err != nil {
		log.Println(err)
		return
	}

	p[peeps.Name] = &Person{Name: peeps.Name, Age: peeps.Age, FavColor: peeps.FavColor}

	jsonValue, _ := json.MarshalIndent(peeps, "", " ")
	fmt.Fprintf(w, "New Person Created\n")
	fmt.Fprintf(w, string(jsonValue))
}

//returnPeople returns everyone in the map
func returnPeople() ([]byte, error) {
	pay = Payload{p}
	return json.MarshalIndent(pay, "", " ")
}

//getEveryone prints out the marshalled json data of everyone in the map
func getEveryone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonValue, err := returnPeople()
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(jsonValue))
	fileHandling(jsonValue)

}

//fileHandling creates a json file and updates it with json data
func fileHandling(jsonValue []byte) {
	f, err := os.Create("people.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.Write(jsonValue)
	if err != nil {
		fmt.Println(err)
		fmt.Println(l)
		f.Close()
		return
	}
}

//findPerson creates a map format and returns the marshalled person you are searching for
func findPerson(found *Person) ([]byte, error) {
	find := make(map[string]*Person)
	find[found.Name] = found
	return json.MarshalIndent(find, "", " ")
}

//getPerson looks for the person you are searching for
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		a := r.URL.Path
		u, err := url.Parse(a)
		if err != nil {
			fmt.Println(err)
		}
		s := strings.Split(u.Path, "/")
		search := s[2]

		for _, item := range p {
			if item.Name == search {
				jsonValue, _ := findPerson(item)
				fmt.Fprintf(w, string(jsonValue))
				return
			}
		}

		fmt.Fprintf(w, "Person Does Not Exist")
	} else {
		fmt.Fprintf(w, "Incorrect Method")
	}
}

func main() {
	http.HandleFunc("/people", peopleRequest)
	http.HandleFunc("/people/", getPerson)
	http.HandleFunc("/people/exit", exitProgram)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//exitProgram kills the program and session
func exitProgram(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exiting Program...")
	os.Exit(1)
}
