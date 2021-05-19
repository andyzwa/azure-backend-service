// API REST EXAMPLE
//
// This is a example over how to create the api from the source.
//
//     Schemes: http
//     Host: localhost:8000
//     Version: 0.1.0
//     basePath: /
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Article - Our struct for all articles
// swagger:model Article
type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Desc    string `json:"color"`
}

// swagger:model Articles
var Articles []Article

func homePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Endpoint Hit: homePage")
	_, err := fmt.Fprintf(w, "Welcome to the HomePage!")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

// swagger:operation GET /articles Articles returnAllArticles
//
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: successful operation
func returnAllArticles(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	err := json.NewEncoder(w).Encode(Articles)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

// swagger:operation GET /article/{id} Articles returnSingleArticle
//
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of article to get
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
//     description: successful operation
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range Articles {
		if article.Id == key {
			err := json.NewEncoder(w).Encode(article)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
		}
	}
}

// swagger:operation POST /article Articles createNewArticle
//
// ---
// parameters:
// - name: article
//   in: body
//   description: article to create or update
//   required: true
//   schema:
//       "$ref": "#/definitions/Article"
// produces:
// - application/json
// responses:
//   '200':
//     description: Article response
//     schema:
//       "$ref": "#/definitions/Article"
func createNewArticle(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: createNewArticle")

	// get the body of our POST request
	// unmarshal this into a new Article struct
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article Article
	err := json.Unmarshal(reqBody, &article)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

	//Append to articles if new
	if article.Id == "" {
		article.Id = strconv.Itoa(len(Articles) + 1)
		Articles = append(Articles, article)
	} else {
		//update
		for index, art := range Articles {
			if art.Id == article.Id {
				Articles[index].Desc = article.Desc
				Articles[index].Title = article.Title
				Articles[index].Content = article.Content
			}
		}
	}

	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}

// swagger:operation DELETE /article/{id} Articles deleteArticle
//
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of article to delete
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
//     description: successful operation
func deleteArticle(_ http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST", "PUT")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	myRouter.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/"))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin", "Content-Type"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Insert the middleware
	log.Fatal(http.ListenAndServe(":8000", c.Handler(myRouter)))

}

func main() {

	//Init some example articles
	Articles = []Article{
		{Id: "1", Title: "Article 1", Content: "Article Content 1", Desc: "#ff2"},
		{Id: "2", Title: "Article 2", Content: "Article Content 2", Desc: "#0bdcab"},
	}

	//Start http listener
	handleRequests()
}
