package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mazen160/go-random"
)

const file string = "routes.db"

const createTable string = `
	CREATE TABLE IF NOT EXISTS Routes(
		url TEXT PRIMARY KEY,
		shortURL TEXT
	)

`

type urlRoute struct {
	longURL  string
	shortURL string
}

func initDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err == nil {
		_, err := db.Exec(createTable)
		if err == nil {
			return db, nil
		} else {
			fmt.Println(err)
			return nil, errors.New("Unable to open database")
		}
	} else {
		return nil, errors.New("Unable to open database")
	}
}

func shortenURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("No URL provided")
	}

	data, err := random.String(9)
	if err != nil {
		return "", errors.New("Unable to generate url")
	}

	return data, nil

}

func insertRoute(db *sql.DB, r urlRoute) (bool, error) {
	_, err := db.Exec(`INSERT INTO Routes (url,shortURL) 
		VALUES(?,?);`, r.longURL, r.shortURL)
	if err == nil {
		return true, err
	}

	return false, err

}

func getRoute(db *sql.DB, s *urlRoute) error {

	rows, err := db.Query("select * from Routes where shortURL = ?", s.shortURL)
	if err == nil {
		for rows.Next() {
			err := rows.Scan(&s.longURL, &s.shortURL)
			if err != nil {
				fmt.Println(err)
				return errors.New("Unable to get redirect root")
			}
		}
		rows.Close()
	} else {
		fmt.Println(err)
	}

	return err
}

func getLongRoute(db *sql.DB, s *urlRoute) error {

	rows, err := db.Query("select * from Routes where url = ?", s.longURL)
	if err == nil {
		for rows.Next() {
			err := rows.Scan(&s.longURL, &s.shortURL)
			if err != nil {
				fmt.Println(err)
				return errors.New("Unable to get redirect root")
			}
		}
		rows.Close()
	} else {
		fmt.Println(err)
	}

	return err

}

func isUniqueEntry(db *sql.DB, url string) bool {

	h := urlRoute{longURL: "", shortURL: ""}

	rows, err := db.Query("select * from Routes where url=?", url)

	if err == nil {

		for rows.Next() {
			err := rows.Scan(&h.longURL, &h.shortURL)
			if err != nil {
				fmt.Println(err)
				return false
			}
		}
		rows.Close()
	} else {
		log.Println(err)
		return false
	}

	return h.longURL != url

}
func main() {

	db, _ := initDBConnection()

	defer db.Close()

	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/index.html")

		if err == nil {
			type data struct {
				Home string
			}

			n := data{"GoURL Shortener"}
			tmpl.Execute(w, n)
		}

		log.Println(r.Method)
	})

	router.HandleFunc("/url/{[a-Z]\\w+}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Called")
		log.Println(r.URL)
		vars := mux.Vars(r)
		short := vars["[a-Z]\\w+"]

		type shortRoute struct {
			shortURL string
		}

		df := urlRoute{longURL: "", shortURL: short}
		err := getRoute(db, &df)

		if err != nil {
			http.Redirect(w, r, df.longURL, http.StatusNotFound)
			return
		}

		if df.longURL != "" {
			w.Header().Set("Content-Type", "")
			http.Redirect(w, r, df.longURL, http.StatusTemporaryRedirect)

		}

	})

	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			data := r.FormValue("url")

			log.Println(data)

			element := ` <textarea name="read-only" readonly> {{.ShortURL}}</textarea>`
			errMsg := `<textarea> name="read-only" Unable to shorten url </textarea>`

			var e error

			tmpl, templateErr := template.New("element").Parse(element)

			if isUniqueEntry(db, data) {

				fmt.Println("Should not be here")

				shrt, urlError := shortenURL(data)

				if urlError != nil {
					fmt.Fprintln(w, errMsg)
					return
				}

				retString := fmt.Sprintf("localhost:9032/url/%v", shrt)

				tmpRoute := urlRoute{longURL: data, shortURL: shrt}
				_, insertErr := insertRoute(db, tmpRoute)
				e = insertErr

				if e != nil || templateErr != nil {
					fmt.Fprintln(w, errMsg)
					return
				} else {
					fmt.Println(e)
				}

				ret := struct {
					ShortURL string
				}{
					retString,
				}

				tmpl.Execute(w, ret)

			} else {

				fmt.Println("Should be here")

				df := urlRoute{longURL: data, shortURL: ""}

				err := getLongRoute(db, &df)

				if err != nil && templateErr != nil {
					fmt.Fprintln(w, errMsg)
					return
				}

				retString := fmt.Sprintf("localhost:9032/url/%v", df.shortURL)

				ret := struct {
					ShortURL string
				}{
					retString,
				}

				tmpl.Execute(w, ret)

			}

			//Return html with shortended url
		}
	})

	http.ListenAndServe(":9032", router)

}