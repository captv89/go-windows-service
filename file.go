package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Files struct {
	Files []File `json:"list"`
}

type File struct {
	Title string `json:"title"`
	Sync  string `json:"sync"`
	Done  string `json:"done"`
}

type conf struct {
	Azvideofolder    string `yaml:"azvideofolder"`
	Localvideofolder string `yaml:"localvideofolder"`
	Azpdffolder      string `yaml:"azpdffolder"`
	Localpdffolder   string `yaml:"localpdffolder"`
	Videojsonpath    string `yaml:"videojsonpath"`
	Pdfjsonpath      string `yaml:"pdfjsonpath"`
	Dbserver         string `yaml:"dbserver"`
	Dbname           string `yaml:"dbname"`
}

func (c *conf) getConfig() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func fileAction(filename string) {
	//read the byte value
	file, err := ioutil.ReadFile(filename)
	errorHandler(err)

	data := Files{}

	byteValue := []byte(file)
	//convert byte to readable data
	err = json.Unmarshal(byteValue, &data)
	errorHandler(err)

	//loop over data
	for i := 0; i < len(data.Files); i++ {
		if filename == videojson {
			fmt.Println("Checking if VIDEO dowload required")
			ans := checkAction(data.Files[i].Title, data.Files[i].Sync, data.Files[i].Done, "video")
			if ans == "yes" {
				data.Files[i].Sync = "no"
				data.Files[i].Done = "yes"
				// Convert golang object back to byte
				// var err error
				byteValue, err = json.Marshal(data)
				errorHandler(err)

				// Write back to file
				err = ioutil.WriteFile(filename, byteValue, 0644)
				errorHandler(err)
			}
		} else {
			fmt.Println("Checking if PDF dowload required")
			ans := checkAction(data.Files[i].Title, data.Files[i].Sync, data.Files[i].Done, "pdf")
			if ans == "yes" {
				data.Files[i].Sync = "no"
				data.Files[i].Done = "yes"
				// Convert golang object back to byte
				// var err error
				byteValue, err = json.Marshal(data)
				errorHandler(err)

				// Write back to file
				err = ioutil.WriteFile(filename, byteValue, 0644)
				errorHandler(err)
			}
		}
	}
}

func checkAction(title string, sync string, done string, kind string) string {
	if sync != "yes" && done == "yes" {
		fmt.Println("Download Not Required For", title)
	} else {
		fmt.Println("Download Required For", title)
		filename := title + ".mp4"

		var e error

		if kind == "video" {
			e = download(azurevideo, filename)
			if e != nil {
				makeEntry(title, "failed", "download", e.Error())
			} else {
				makeEntry(title, "success", "download", "")
			}

			e = moveFile(videofolder, filename)
			if e != nil {
				makeEntry(title, "failed", "move", e.Error())
			} else {
				makeEntry(title, "success", "move", "")
				return "yes"
			}
		} else {
			e = download(azurepdf, filename)
			if e != nil {
				makeEntry(title, "failed", "download", e.Error())
			} else {
				makeEntry(title, "success", "download", "")
			}

			e = moveFile(pdffolder, filename)
			if e != nil {
				makeEntry(title, "failed", "move", e.Error())
			} else {
				makeEntry(title, "success", "move", "")
				return "yes"
			}
		}
	}
	return "no"
}

func moveFile(foldername string, filename string) (e error) {
	newpath := foldername + filename
	err := os.Rename(filename, newpath)
	if err != nil {
		fmt.Println(err)
		return e
	}
	return
}
