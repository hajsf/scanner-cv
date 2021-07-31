package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zserge/lorca"
)

//go:embed static
var staticFiles embed.FS

func main() {
	go func() {
		http.HandleFunc("/favicon.ico", func(rw http.ResponseWriter, r *http.Request) {})
		// http.FS can be used to create a http Filesystem
		var staticFS = http.FS(staticFiles)
		fs := http.FileServer(staticFS) // embeded static files
		// Serve static files, to be embedded in the binary
		http.Handle("/static/", fs)

		//	www := http.FileServer(http.Dir("/files/")) // side static files
		// Serve public files, to be beside binary
		http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./files"))))

		http.HandleFunc("/getSkills", getSkills)
		log.Println("Listening on :3000...")
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Start UI
	ui, err := lorca.New("http://localhost:3000/static/", "", 1200, 800)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer ui.Close()

	// Bind Go function to be available in JS. Go function may be long-running and
	// blocking - in JS it's represented with a Promise.
	ui.Bind("add", func(a, b int) int { return a + b })

	// Call JS function from Go. Functions may be asynchronous, i.e. return promises
	n := ui.Eval(`Math.random()`).Float()
	fmt.Println(n)

	// Call JS that calls Go and so on and so on...
	m := ui.Eval(`add(2, 3)`).Int()
	fmt.Println(m)

	// Wait for the browser window to be closed
	<-ui.Done()
}

/*
To return JSON to the client instead of text
	data := SomeStruct{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
*/
func getSkills(w http.ResponseWriter, r *http.Request) {
	// define custom type
	type Input struct {
		Path   string `json:"path"`
		Skills string `json:"skills"`
	}

	// define a var
	var input Input

	// decode input or return error
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Decode error! please check your JSON formating.")
		return
	}

	// print user inputs
	fmt.Fprintf(w, "Inputed name: %s", input.Path)
}
