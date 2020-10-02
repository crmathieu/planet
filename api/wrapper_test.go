package api

import (
	"net/http"
	"testing"
//    "github.com/crmathieu/planet/data"
//    "encoding/json"
    //"fmt"
    "bytes"
    "io/ioutil"
)
/*
func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestUsersFamily(t *testing.T) {

    req, _ := http.NewRequest("GET", "/products", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}
*/

func DELETE(t *testing.T, url string, dummy string) (int, string, error){

    req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte(dummy)))
    if err != nil {
        t.Fatal(err)
    }
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatal(err)
    }

    defer resp.Body.Close()
    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    // Convert response body to string
    bodyString := string(bodyBytes)
    return resp.StatusCode, bodyString, err
}

func PUT(t *testing.T, url string, payload string) (int, string, error){
    req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(payload)))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json; charset=utf-8")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatal(err)
    }

    defer resp.Body.Close()
    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    // Convert response body to string
    bodyString := string(bodyBytes)
    return resp.StatusCode, bodyString, err

    // Convert response body to Todo struct
//   var todoStruct Todo
//   json.Unmarshal(bodyBytes, &todoStruct)
//   fmt.Printf("API Response as struct:\n%+v\n", todoStruct)
}

func POST(t *testing.T, url string, payload string) (int, string, error){
    resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(payload)))
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()
    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    // Convert response body to string
    bodyString := string(bodyBytes)
    return resp.StatusCode, bodyString, err
}

func GET(t *testing.T, url string, dummy string) (int, string, error){
    resp, err := http.Get(url)
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()
    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    // Convert response body to string
    bodyString := string(bodyBytes)
    return resp.StatusCode, bodyString, err
}

var httpVerbs = map[string]func(*testing.T, string, string) (int, string, error){
    "PUT": PUT,
    "DELETE": DELETE,
    "POST": POST,
    "GET": GET,
}

type TESTdata struct {
        comment string
        method  string
        url     string
        requestBody    string
        expected int
        checkBody bool
        checkBodyFunction func(string, string) bool
        results string
}

func TestApi(t *testing.T) {

    var testapp = []TESTdata{
        {
            "get a non-existent user",
            "GET",
            "http://localhost/users/raoulp",
            "",
            404,
            false,
            nil,
            "",
        },
        {
            "Add a user",
            "POST",
            "http://localhost/users",
            `{"first_name": "raoul", "last_name": "popov", "userid": "raoulp", "groups": ["newgrp", "users"]}`,
            200,
            false,
            nil,
            "",
        },
        {
            "delete a non-existing user",
            "DELETE",
            "http://localhost/users/samp",
            "",
            404,
            false,
            nil,
            "",
        },
        {
            "delete an existing user",
            "DELETE",
            "http://localhost/users/raoulp",
            "",
            200,
            false,
            nil,
            "",
        },
        {
            "get an existing user",
            "GET",
            "http://localhost/users/joes",
            "",
            200,
            false,
            nil,
            "",
        },
        {  
            "update a non-existent user",
            "PUT",
            "http://localhost/users/raoulp",
            `{"first_name": "raoul", "last_name": "popov", "userid": "raoulp", "groups": ["newgrp", "users"]}`,
            404,
            false,
            nil,
            "",
        },  
        {
            "update an existing user",
            "PUT",
            "http://localhost/users/joes",
            `{"first_name": "joseph", "last_name": "smith", "userid": "joes", "groups": ["admin", "newgrp", "users"]}`,
            200,
            false,
            nil,
            "",
        },
    }

    for _, v := range testapp {
        t.Log(v.comment)
        statusCode, body, err := httpVerbs[v.method](t, v.url, v.requestBody)
        if err != nil {
            t.Fatal(err)
        }
        if statusCode != v.expected {
            t.Errorf("expected %v, got %d\nBody contains: %v", v.expected, statusCode, body)
        } else {
            t.Logf("status %d as expected\n", statusCode)
        }

        if v.checkBody {
            v.checkBodyFunction(body, v.results)
        }

        t.Log("\n")
        //var dst struct{ Salutation string }
        //if v.checkBody {
        //    if err := v.verify(body); err != nil {
        //        t.Fatal(err);
        //    }
       // }
/*        var user data.USER 
        ok := json.Unmarshal([]byte(body), &user)
        if ok != nil {
            t.Fatal(ok)
        }
        fmt.Println(user)
        if v.results != body {

        }
*/
    }
/*    resp, err := http.Get("http://localhost/users/joes")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
    //var dst struct{ Salutation string }
    var user data.USER 
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
    }
    fmt.Println(user)*/
//	if dst.Salutation != "Hello Frank!" {
//		t.Fatalf("expected 'Hello Frank!', got '%s'", dst.Salutation)
//	}
}