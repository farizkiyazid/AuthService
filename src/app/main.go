package main
 
import (
    "fmt"
    "log"
	"net/http"
	"html/template"
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
    "mockDB"
    "encoding/json"
    "bytes"
)



func hello(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
 
    switch r.Method {
    case "GET":     
         http.ServeFile(w, r, "form.html")
    case "POST":
        // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
        if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        // fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
        user := r.FormValue("user")
        password := r.FormValue("password")
        if user==mockDB.User{
            if password==mockDB.Pass{
                token := GenerateToken(user+password)
                myvar := map[string]interface{}{"token": token}
                MakeRequest(user, token)
                outputHTML(w, "hasil.html", myvar)
            }else{
                fmt.Fprintf(w, "Can't find matching password")
            }
        }else{
            fmt.Fprintf(w, "Can't find matching user")
        }
    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
    t, err := template.ParseFiles(filename)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if err := t.Execute(w, data); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}

// GenerateToken returns a unique token based on the provided email string
func GenerateToken(email string) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Hash to store:", string(hash))

    hasher := md5.New()
    hasher.Write(hash)
    return hex.EncodeToString(hasher.Sum(nil))
}

func MakeRequest(user string, token string) {
    //change url for bookservice
    urlNya := "http://localhost:8080/test"

	message := map[string]interface{}{
		"user": user,
		"token":  token,
    }

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
    }
    log.Println(message)
    log.Println(bytesRepresentation)

	resp, err := http.Post(urlNya, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
}
 
func main() {
    http.HandleFunc("/", hello)
 
    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
