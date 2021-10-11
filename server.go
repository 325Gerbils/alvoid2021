package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

)

/*
`survey data.json`
{
    "Name": "Alvoid Survey",
    "Questions": [
        "Can you bring a tent?",
        "Can you bring water?",
        "Can you find a light?"
    ],
    "Responders": {
        "Alvoid 1": {
            "Name": "Alvoid 1",
            "Answers": {
                "Can you bring a tent?": "I can bring a tent.",
                "Can you bring water?": "I can bring water.",
                "Can you find a light?": "I can find a light."
            }
        },
        "Alvoid 2": {
            "Name": "Alvoid 2",
            "Answers": {
                "Can you bring a tent?": "I can bring a tent.",
                "Can you bring water?": "I can bring water.",
                "Can you find a light?": "I can find a light."
            }
        }
    }
}
*/

type Survey struct {
	Name       string
	Questions  []string
	Responders map[string]*Responder
}

type Responder struct {
	Name    string
	Answers map[string]string
}

func (s *Survey) AddQuestion(question string) {
	s.Questions = append(s.Questions, question)
}

func (s *Survey) AddResponses(responderName string, answers map[string]string) {
	s.Responders[responderName] = &Responder{
		Name:    responderName,
		Answers: answers,
	}
}

func backupSurvey(survey *Survey) {
	f, err := os.Create("db.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(survey)
}

func main() {
	var survey *Survey = &Survey{}

	// load survey id backup exists in .db, otherwise initialize w questions etc
	if file, err := os.Open("./db.json"); err == nil {
		log.Println("Loading survey from backup")
		decoder := json.NewDecoder(file)
		err := decoder.Decode(survey)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	} else {
		survey.Name = "Alvoid Survey"

		survey.Questions = []string{"Are you attending Alvoid?",
			"What is your name?",
			"How many other people are you bringing?",
			"Do you want to/ are you going to perform?",
			"What is your artist's name?",
			"Where are you from?",
			"How would you describe your music?",
			"Please link your socials and streaming platforms",
			"Can you provide (or need) transportation?",
			"What infrastructure/useful items can you bring?"}
		survey.Responders = make(map[string]*Responder)

		log.Println("Initialized survey")
		backupSurvey(survey)
		log.Println("Saved survey to backup")
	}

	http.HandleFunc("/respond", func(w http.ResponseWriter, r *http.Request) {
		// get POSTed JSON and unmarshal it into a Responder, then add it to survey
		var responder *Responder
		err := json.NewDecoder(r.Body).Decode(&responder)
		log.Println("Got response:", responder)
		if err != nil {
			log.Println(err)
			return
		}
		survey.AddResponses(responder.Name, responder.Answers)
		// write contents of survey to .db
		backupSurvey(survey)
		log.Println("Saved survey to backup")
		fmt.Fprintf(w, "Responder %s added to survey", responder.Name)
	})

	// write survey contents to /view
	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving survey contents...")
		fmt.Fprintf(w, "<html><head><style>td{border:1px solid lightgrey;padding:5px;}</style></head><body>")
		fmt.Fprintf(w, "<h1>"+survey.Name+"</h1><br>")
		for _, responder := range survey.Responders {
			fmt.Fprintf(w, "<table><tr><th colspan=\"2\">"+responder.Name+"</th></tr>")
			for question, answer := range responder.Answers {
				fmt.Fprintf(w, "<tr><td>"+question+"</td>")
				fmt.Fprintf(w, "<td>"+answer+"</td></tr>")
			}
			fmt.Fprintf(w, "</table><br>")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})
	
	http.HandleFunc("/images/alvord.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord.jpg")
	})
	http.HandleFunc("/images/alvord_bikes.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord_bikes.jpg")
	})
	http.HandleFunc("/images/alvord_bushes.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord_bushes.jpg")
	})
	http.HandleFunc("/images/alvord_sunset.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord_sunset.jpg")
	})
	http.HandleFunc("/images/alvord_sunset2.jpeg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord_sunset2.jpeg")
	})
	http.HandleFunc("/images/alvord3.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord3.jpg")
	})
	http.HandleFunc("/images/alvord4.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "images/alvord4.jpg")
	})

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
