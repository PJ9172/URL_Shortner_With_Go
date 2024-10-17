package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type UrlStruct struct {
	Id          string    `json:"id"`
	OrignalUrl  string    `json:"orignalurl"`
	ShortUrl    string    `json:"shorturl"`
	CurrentDate time.Time `json:"currentdate"`
}

var urlDB = make(map[string]UrlStruct)

func main() {

	/*fmt.Print("Enter the url : ")
	reader := bufio.NewReader(os.Stdin)
	URL, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error to reading url", err)
		return
	}

	URL = strings.TrimSpace(URL)
	// fmt.Println(URL)

	shorturl := storeInStruct(URL)
	// fmt.Println(shorturl)

	newurl := "http://localhost:3000/urlshorter/" + shorturl
	fmt.Println("\nSearch the following URL : ", newurl)*/

	//Register the handler function to all request
	http.HandleFunc("/", rootUrlHandler)
	http.HandleFunc("/submit", shortUrl)
	http.HandleFunc("/urlshorter/", redirectUrl)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	//Server starting on 3000 port
	fmt.Println("Server starting on 3000 port!!!")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Print("Error in server starting", err)
		return
	}
}

func storeInStruct(orignalurl string) string {
	shorturl := generateShortUrl(orignalurl)
	id := shorturl
	urlDB[id] = UrlStruct{
		Id:          id,
		OrignalUrl:  orignalurl,
		ShortUrl:    shorturl,
		CurrentDate: time.Now(),
	}
	return shorturl
}

func generateShortUrl(orignalurl string) string {
	hasher := md5.New()
	hasher.Write([]byte(orignalurl))
	hasherSlice := hasher.Sum(nil)
	hash := hex.EncodeToString(hasherSlice)
	shorturl := hash[:8]
	return shorturl
}

func getOrignalUrl(id string) (UrlStruct, error) {
	urlstruct, ok := urlDB[id]
	if !ok {
		return UrlStruct{}, fmt.Errorf("url not found")
	}
	return urlstruct, nil
}

func redirectUrl(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/urlshorter/"):]

	// Fetch original URL using the shortened ID
	urlstruct, err := getOrignalUrl(id)
	if err != nil {
		// If the URL is not found, return a 404 error
		http.Error(res, "Invalid request!! URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(res, req, urlstruct.OrignalUrl, http.StatusFound)
}

func shortUrl(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello Prajwal...")
	if req.Method == "POST" {
		req.ParseForm()
		URL := req.FormValue("URL")
		shorturl := storeInStruct(URL)
		newUrl := "http://localhost:3000/urlshorter/" + shorturl
		fmt.Fprint(res, "Your URL is Successfully Shorted !!!\nSearch it : ", newUrl)
		fmt.Println(newUrl)
	}
}

func rootUrlHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}
