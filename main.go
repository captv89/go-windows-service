package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var azurevideo, videofolder, azurepdf, pdffolder, videojson, pdfjson, dbserver, dbname string // Config Declaration

var vdojsonmod fs.FileInfo
var pdfjsonmod fs.FileInfo

func init() {
	// Load Env Variables
	godotenv.Load()
	// Load Config
	var c conf
	c.getConfig()
	azurevideo = c.Azvideofolder
	videofolder = c.Localvideofolder
	azurepdf = c.Azpdffolder
	pdffolder = c.Localpdffolder
	videojson = c.Videojsonpath
	pdfjson = c.Pdfjsonpath
	dbserver = c.Dbserver
	dbname = c.Dbname

	// Open DB Connection
	var e error
	db, e = sql.Open(dbserver, dbname)
	errorHandler(e)

	createDB()

	// regiter initial time

	vdojsonmod, e = os.Stat(videojson)
	errorHandler(e)
	pdfjsonmod, e = os.Stat(pdfjson)
	errorHandler(e)

	for {
		curvdoStat, err := os.Stat(videojson)
		errorHandler(err)
		curpdfStat, err := os.Stat(pdfjson)
		errorHandler(err)
		rand.Seed(time.Now().UnixNano())
		if vdojsonmod.ModTime() != curvdoStat.ModTime() || pdfjsonmod.ModTime() != curpdfStat.ModTime() {
			fmt.Println("File Modified.. Trying to Download updated files..")
			vdojsonmod = curvdoStat
			pdfjsonmod = curpdfStat
			main()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		} else {
			fmt.Println("Files not modified, going to sleep..")
			time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
			fmt.Println("Sleep Over.. Let me check again..")
		}

	}
}

func main() {
	fileAction(videojson)
	fileAction(pdfjson)
	printTable()
}

func errorHandler(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
