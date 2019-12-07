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
        // fmt.Fprintf(w, "Name = %s\n", name)
		// fmt.Fprintf(w, "Password = %s\n", password)
        if user==mockDB.User{
            if password==mockDB.Pass{
                token := GenerateToken(user+password)
                myvar := map[string]interface{}{"token": token}
                outputHTML(w, "hasil.html", myvar)
            }
        }else{
            fmt.Fprintf(w, "Can't find matching user/password")
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
 
func main() {
    http.HandleFunc("/", hello)
 
    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
