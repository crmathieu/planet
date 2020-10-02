package api

import (
	"net/http"
	"testing"
    "bytes"
    "io/ioutil"
    "fmt"
)


var httpVerbs = map[string]func(*testing.T, string, string) (int, string, error){
    "PUT":      usePUT,
    "DELETE":   useDELETE,
    "POST":     usePOST,
    "GET":      useGET,
}

type TESTdata struct {
        comment string
        method  string
        url     string
        requestBody    string
        expected int
        checkBody bool
        checkBodyFunction func(*testing.T, string, string) error
        results string
}

func TestApi(t *testing.T) {

    var testapp = []TESTdata{
        {
            "\nGet a non-existent user",
            "GET",
            "http://localhost/users/raoulp",
            "",
            404,
            false,
            nil,
            "",
        },
        {
            "\nAdd a user",
            "POST",
            "http://localhost/users",
            `{"first_name": "raoul", "last_name": "popov", "userid": "raoulp", "groups": ["newgrp", "users"]}`,
            200,
            false,
            nil,
            "",
        },
        {
            "\nDelete a non-existing user",
            "DELETE",
            "http://localhost/users/samp",
            "",
            404,
            false,
            nil,
            "",
        },
        {
            "\nDelete an existing user",
            "DELETE",
            "http://localhost/users/raoulp",
            "",
            200,
            false,
            nil,
            "",
        },
        {
            "\nGet an existing user",
            "GET",
            "http://localhost/users/joes",
            "",
            200,
            false,
            nil,
            "",
        },
        {  
            "\nUpdate a non-existent user",
            "PUT",
            "http://localhost/users/raoulp",
            `{"first_name": "raoul", "last_name": "popov", "userid": "raoulp", "groups": ["newgrp", "users"]}`,
            404,
            false,
            nil,
            "",
        },  
        {
            "\nUpdate an existing user",
            "PUT",
            "http://localhost/users/joes",
            `{"first_name": "joseph", "last_name": "smith", "userid": "joes", "groups": ["admin", "newgrp", "users"]}`,
            200,
            false,
            nil,
            "",
        },
        {
            "\nGet group members",
            "GET",
            "http://localhost/groups/users",
            "",
            200,
            true,
            displayResult,
            `["phils","johnf","joes"]`,
        },
        {
            "\nAdd an empty group",
            "POST",
            "http://localhost/groups",
            `{"group_name":"anothergrp"}`,
            200,
            false,
            nil,
            "",
        },
        {
            "\nDelete a group",
            "DELETE",
            "http://localhost/groups/anothergrp",
            "",
            200,
            false,
            nil,
            "",
        },
        {
            "\nUpdate a group's member list",
            "PUT",
            "http://localhost/groups/users",
            `["joes","johnf"]`,
            200,
            false,
            nil,
            "",
        },

    }

    for _, v := range testapp {
       
        testname := fmt.Sprintf("%v", v.comment)
        t.Run(testname, func(t *testing.T) {
        
            statusCode, body, err := httpVerbs[v.method](t, v.url, v.requestBody)
            if err != nil {
                t.Fatal(err)
            }
            if statusCode != v.expected {
                t.Errorf("expected %v, got %d\nBody contains: %v", v.expected, statusCode, body)
            }

            if v.checkBody {
                e := v.checkBodyFunction(t, body, v.results)
                if e != nil {
                    t.Log(e.Error())
                }
            }
        })
    }
}

func useDELETE(t *testing.T, url string, dummy string) (int, string, error){

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

func usePUT(t *testing.T, url string, payload string) (int, string, error){
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
}

func usePOST(t *testing.T, url string, payload string) (int, string, error){
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

func useGET(t *testing.T, url string, dummy string) (int, string, error){
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

func displayResult(t *testing.T, body string, results string) error {
    t.Logf("\nresults: ==> %v\n\n", body)

    // this function could be customized for specific endpoints to verify that the 
    // payload returned is the one expected.
    return nil
}


