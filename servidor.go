package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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

func main() {
	http.HandleFunc("/mujeres", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Response Ok")
			return

		case http.MethodGet:
			totalResults := 25000
			resultsPerPage := 1000
			numRequests := totalResults / resultsPerPage

			allUsers := make([]User, 0)
			usersCh := make(chan []User)
			wg := sync.WaitGroup{}
			countMutex := sync.Mutex{}
			resultCount := 0

			for i := 0; i < numRequests; i += 5 {
				wg.Add(5)
				for j := 0; j < 5; j++ {
					go func() {
						defer wg.Done()
						url := fmt.Sprintf("https://randomuser.me/api/?results=%d", resultsPerPage)
						response, err := http.Get(url)
						if err != nil {
							log.Printf("Error fetching data: %s", err.Error())
							return
						}
						defer response.Body.Close()

						body, err := ioutil.ReadAll(response.Body)
						if err != nil {
							log.Printf("Error reading response body: %s", err.Error())
							return
						}

						var randomUserResp Mujeres
						err = json.Unmarshal(body, &randomUserResp)
						if err != nil {
							log.Printf("Error parsing JSON response: %s", err.Error())
							time.Sleep(10 * time.Millisecond) // Wait for 10 milliseconds before continuing
						}

						// Count the number of results
						countMutex.Lock()
						resultCount += len(randomUserResp.Results)
						countMutex.Unlock()

						usersCh <- randomUserResp.Results
					}()
				}

				time.Sleep(2 * time.Second) // Sleep for 2 seconds before the next set of requests
			}

			go func() {
				wg.Wait()
				close(usersCh)
			}()

			for users := range usersCh {
				allUsers = append(allUsers, users...)
			}

			responseData, err := json.Marshal(allUsers)
			if err != nil {
				http.Error(w, "Error marshaling response data", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(responseData)

			fmt.Printf("Total Results: %d\n", resultCount)
		}
	})

	log.Println("Server listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
