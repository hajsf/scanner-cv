package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func getSkills(w http.ResponseWriter, r *http.Request) {
	// define custom types
	type Input struct {
		Path   string `json:"path"`
		Skills string `json:"skills"`
	}
	type Output struct {
		MSG    string   `json:"msg"`
		File   string   `json:"file"`
		Skills []string `json:"skills"`
	}

	// define vars
	var input Input
	var output Output

	// decode input or return error
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Decode error! please check your JSON formating.")
		return
	}

	http.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(input.Path))))

	tmpDir, exists := os.LookupEnv("cvTemp")
	fmt.Println("cvTemp: ", exists)
	fmt.Println("tmpDir: ", tmpDir)
	if !exists {
		tmpDir = os.TempDir() + "\\cvscanner"
		err = exec.Command(`SETX`, `cvTemp`, tmpDir).Run()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		fmt.Println("tmpDir: ", tmpDir)
		//	err = os.Setenv("keyTemp", tmpDir)
		if err != nil {
			panic(err)
		}
	}

	var files []string
	var paths []string

	root := input.Path
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.ToLower(filepath.Ext(path)) == ".pdf" {
			f := strings.TrimSuffix(info.Name(), filepath.Ext(path))
			files = append(files, f)
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for i, file := range files {
		fmt.Println(file + ".pdf")
		str, err := ExtractTexts(file, paths[i])
		if err != nil {
			panic(err)
		}
		re := regexp.MustCompile(`(?i)` + input.Skills)
		matches := re.FindAllString(str, -1)
		uniqueSlice := unique(matches)
		p := strings.Split(input.Skills, "|")
		fmt.Println("skills found:", len(uniqueSlice), "/", len(p))
		fmt.Println("'" + strings.Join(uniqueSlice, `', '`) + `'`)
		output = Output{MSG: "ok", File: file, Skills: uniqueSlice}
		data, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		//	w.WriteHeader(http.StatusCreated)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		//	json.NewEncoder(w).Encode(output)
	}
	// print user inputs
	//fmt.Fprintf(w, "Inputed name: %s", input.Path)
}
