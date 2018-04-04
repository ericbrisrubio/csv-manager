package main

import (
	"log"
	"os"
	"github.com/urfave/cli"
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"encoding/json"
	"github.com/satori/go.uuid"
	_ "io"
	"io"
	"fmt"
	"strings"
)

var jsonFile ConfigFile

func main() {
	app := cli.NewApp()
	app.Name = "Csv-Manager"
	app.Usage = "Get info from a csv and modify its headers"
	app.Version="0.1.0"
	var filePath string
	var configPath string

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "filepath, f",
			Value:       "",
			Usage:       "the CSV file path",
			Destination: &filePath,
		},
		cli.StringFlag{
			Name:        "configp, c",
			Value:       "",
			Usage:       "the config file path",
			Destination: &configPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "headers",
			Aliases: []string{"d"},
			Usage:   "Modify headers for a file",
			Action:  func(c *cli.Context) error {
				if configPath != "" {
					if filePath != ""{
						ModifyHeaders(filePath)
					} else {
						log.Println("`filepath, f` flag has to be provided")
					}
				} else {
					log.Println("`configp, c` flag has to be provided")
				}
				return nil
			},
		},
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "info about the file",
			Action:  func(c *cli.Context) error {
				//loadConfig(configPath)
				//log.Println(jsonFile)
				if filePath != ""{
					log.Println("Getting the info for your csv... (it could take a few seconds)")
					//Create the CsvModifier instance
					csvModifier := CsvModifier{filePath, uuid.NewV4().String() + ".csv", nil}
					//only modify the 1st line
					linesTotal := countCsvLines(csvModifier.CsvToModifyPath)
					log.Printf("Total rows: %d", linesTotal)
				} else {
					log.Println("`filepath, f` flag has to be provided")
				}


				return nil
			},
		},
	}
	//sort.Sort(cli.Command(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

/**
Modify the headers of the csv and returns a new csv with the headers changed
 */
func ModifyHeaders(filepath string) {
	loadConfig("config.json")
	//log.Println(jsonFile)
	log.Println("Modifying the headers for your csv!")
	//Create the CsvModifier instance
	csvModifier := CsvModifier{filepath, uuid.NewV4().String() + ".csv", nil}
	//only modify the 1st line
	fileFrom, _ := os.Open(csvModifier.CsvToModifyPath)
	linesTotal := countCsvLines(csvModifier.CsvToModifyPath)
	linesProcessed := 0
	log.Printf("Total lines: %d", linesTotal)
	//defer
	reader := csv.NewReader(bufio.NewReader(fileFrom))
	line, _ := reader.Read()
	newLine := csvModifier.changeNames(line)
	//log.Println(line)
	fileNew := csvModifier.createFile()
	w := csv.NewWriter(fileNew)
	w.Write(newLine)
	linesProcessed++
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	//write the rest of the lines
	for {
		//w1 := csv.NewWriter(fileNew)
		record, err2 := reader.Read()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			log.Fatal(err2)
		}
		//log.Println(record)
		errW := w.Write(record)
		if errW != nil {
			log.Fatal(errW)
		}
		w.Flush()
		linesProcessed++
		fmt.Println(calculatePercentProcessed(linesProcessed, linesTotal))
	}
	fileFrom.Close()
	fileNew.Close()
	log.Printf("%d lines processed", linesProcessed)
}

/**
Calculate the percentage of the lines being processed
 */
func calculatePercentProcessed(processed int, total int) int{
	return int(processed*100/total)
}

/**
This function counts and returns the amount of rows in a csv file
@params filepath -> the path of the file to be used by the function
 */
func countCsvLines(filepath string) int{
	file, _ := os.Open(filepath)
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	file.Close()
	return lineCount
}

func loadConfig(path string) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("%s", err)
	}
	json.Unmarshal(raw, &jsonFile)
}

type ConfigFile struct {
	OldNewNames map[string]string `json:"old_new_names"`
	NameChanger string `json:"name_changer_instance"`
}

type CsvModifier struct {
	CsvToModifyPath string
	ResultCsvPath   string
	File            *os.File
}

/**
Create a new file with a random name uuid
 */
func (cr *CsvModifier) createFile() (*os.File) {
	filename := cr.ResultCsvPath
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file %s.csv", filename)
	} else {
		log.Printf("%s file has been created", filename)
	}
	return file
}

/**
Changes values for a row, in this case it will be used to changed the headers based on the config file
 */
func (cr *CsvModifier) changeNames(vals []string) ([]string) {
	changerName := jsonFile.NameChanger
	changerInstance := FactoryNameChanger[changerName]
	instanceOFChanger := changerInstance()
	return instanceOFChanger.ModifyList(vals)
}



/**
Implement this interface those who need to get a different behavior for the headers values
 */
//////////////////////////////////////// ICHANGER INTERFACE ---------------------------------------


type fn func()IChanger

var FactoryNameChanger = map[string]fn{
	"Upper": func() IChanger{return UpperChanger{}},
}

type IChanger interface{
	ModifyList([]string)[]string
}

/**
This struct is going to convert the headers to upperCase from the values they already have
 */
type UpperChanger struct{
}

//This lines guaranties that UpperChanger is going to implement IChanger interface
var _ IChanger = (*UpperChanger)(nil)

func (uc UpperChanger) ModifyList(elList []string) []string {
	for i, e:= range elList {
		elList[i] = strings.ToUpper(e)
	}
	return elList
}


