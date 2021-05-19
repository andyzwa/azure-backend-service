package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHomePage(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	homePage(w, req)

	expected := "Welcome to the HomePage!"
	actual := w.Body.String()
	if !strings.HasPrefix(actual, expected) {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

func TestGetArticles(t *testing.T) {

	//set test articles
	Articles = []Article{
		{Id: "1", Title: "Article 1", Content: "Article Content 1", Desc: "#ff2"},
		{Id: "2", Title: "Article 2", Content: "Article Content 2", Desc: "#0bdcab"},
	}

	//call func
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	returnAllArticles(w, req)

	//assert
	actual, _ := ioutil.ReadAll(w.Body)
	var aAct []Article
	err := json.Unmarshal(actual, &aAct)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if !reflect.DeepEqual(Articles, aAct) {
		t.Fatalf("Expected %s but got %s", Articles, aAct)
	}
}
