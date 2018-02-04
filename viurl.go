package viurl

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("x...")
	r.ParseForm()
	// fmt.Println(r.Form)
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("schema", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {

	// 	fmt.Println("key: ", k)
	// 	fmt.Println("value:", strings.Join(v, ""))

	// }
	// abort
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("html/index.html")
		t.Execute(w, nil)
	} else {
		ip := r.Form["ipAddr"]
		port := r.Form["port"]
		fmt.Fprintf(w, strings.Join(ip, "")+"\r\n"+strings.Join(port, ""))
	}

	// fmt.Fprintf(w, "hello you !\r\n")
}

// func main() {
func viurl() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
