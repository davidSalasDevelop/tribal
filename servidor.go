package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Mujeres struct {
	Results []User `json:"results"`
}

type User struct {
	Gender string `json:"email"`
	Name   struct {
		First string `json:"first"`
	} `json:"name"`
	Uuid struct {
		Uuid string `json:"uuid"`
	} `json:"login"`
}

// func handleRequest(w http.ResponseWriter, r *http.Request) {

// }

func main() {

	// Define HTTP routes and handlers
	http.HandleFunc("/mujeres", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Resonse Ok")

			return

		case http.MethodGet:
			//REtrieve the las inserted location

			totalResults := 25000
			resultsPerPage := 1000
			numRequests := totalResults / resultsPerPage

			var allUsers []User

			for i := 0; i < numRequests; i++ {
				url := fmt.Sprintf("https://randomuser.me/api/?results=%d", resultsPerPage)
				response, err := http.Get(url)
				if err != nil {
					http.Error(w, "Error fetching data", http.StatusInternalServerError)
					return
				}
				defer response.Body.Close()

				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					http.Error(w, "Error reading response body", http.StatusInternalServerError)
					return
				}

				var randomUserResp Mujeres
				err = json.Unmarshal(body, &randomUserResp)
				if err != nil {
					http.Error(w, "Error parsing JSON response", http.StatusInternalServerError)
					return
				}

				allUsers = append(allUsers, randomUserResp.Results...)
			}

			responseData, err := json.Marshal(allUsers)
			if err != nil {
				http.Error(w, "Error marshaling response data", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(responseData)
		}
	})

	log.Println("Server listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
