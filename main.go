package main 

import (
	"flag"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"os"
	"encoding/csv"
	"strings"
	"strconv"
)
type database struct{
	host string
	username string
	password string
	port int64
	database string
}
type patientRecord struct{
	name string
	age int
	city string
	disease string
}
var (
	isDatabaseIniFile bool
 	databaseIniFile string 
	db database
	dataFile string
)
const (
	section ="postgresql"
	hostOption = "host"
	portOption = "port"
	databaseOption = "database"
	passwordOption = "password"
	iniFormat = "\n\r[postgresql] # This Line Should be As it is with square Brackets\n" +
	"host=localhost # Your Database Host Address\n" +
	"database=dbName # Your Database Name\n" +
	"user=dbUser # Your Database User\n"+
	"password=yourPassword # Your Database Password\n"
)
func readCommandLineArgs(){
	flag.BoolVar(&isDatabaseIniFile,"isDbIniFile",true, "Is Database Initializer File Available?")
	flag.StringVar(&databaseIniFile, "dbIniFile", "database.ini", "Database Initializer File Name." )
	flag.StringVar(&db.host, "host","localhost","The Host at Which database is Hosted.")
	flag.StringVar(&db.username,"username","postgres","Username of The Database.")
	flag.StringVar(&db.database,"db","","Database Name.")
	flag.StringVar(&db.password,"password","", "Password of the Database.")
	flag.StringVar(&dataFile,"dataFile","patients.csv", "DataFile from which you want to read the Data.")
	flag.Int64Var(&db.port,"port",5432,"Port at which Database is running." )
	flag.Parse()
}
func printFileIssue(message string){
	fmt.Println(message)
	fmt.Println("Please Specify",databaseIniFile, "in the following format")
	fmt.Println(string(iniFormat))
}
func readIniFile(parser *configparser.ConfigParser) (ok error){
	options:= make([]string,0)
	options,ok=parser.Options(section)
	if ok!= nil{
		return 
	}
	for _,option:= range options {
		switch option {
		case hostOption:
			db.host, ok = parser.Get(section,option)
			if ok!=nil {
				return
			}
		case passwordOption:
			db.password, ok= parser.Get(section,option)
			if ok!=nil {
				return
			}
		case databaseOption:
			db.database, ok = parser.Get(section,option)
			if ok!=nil {
				return
			}
		case portOption:
			db.port, ok = parser.GetInt64(section,option)
			if ok!=nil {
				return
			}
		}
	}
	return
}
func initDBObjectFromIniFile() (ok error){
	var iniFile *configparser.ConfigParser
	iniFile,ok=configparser.NewConfigParserFromFile(databaseIniFile)
	if ok!=nil {
		return 
	}
	if !iniFile.HasSection(section) {
		return
	}
	ok=readIniFile(iniFile)
	if ok!=nil{
		return
	}
	return
}
const (
	nameColumn="name"
	ageColumn="age"
	cityColumn="city"
	diseaseColumn="disease"
)

func main(){
	
	readCommandLineArgs()

	var columnsIndex map[string] int = map[string] int{
		nameColumn:0,
		ageColumn:0,
		cityColumn:0,
		diseaseColumn:0,
	}
	var columnFound map[string] bool = map[string] bool{
		nameColumn:false,
		ageColumn:false,
		cityColumn:false,
		diseaseColumn:false,
	}
	
	
	if ok:=initDBObjectFromIniFile();ok!=nil {
		printFileIssue(ok.Error())
		os.Exit(-1)
	}
	fileReader,ok:=os.Open(dataFile)
	if ok!=nil {
		fmt.Println("No datafile found:", dataFile)
		os.Exit(-1)
	}
	var header []string
	csvReader:=csv.NewReader(fileReader)
	header,ok=csvReader.Read()
	if ok!=nil {
		fmt.Println(ok)
		os.Exit(-1)
	}

	for index, colName:= range header {
		switch strings.ToLower(colName){
		case nameColumn:
			columnFound[colName]=true
			columnsIndex[colName]=index
		case ageColumn:
			columnFound[colName]=true
			columnsIndex[colName]=index
		case cityColumn:
			columnFound[colName]=true
			columnsIndex[colName]=index
		case diseaseColumn:
			columnFound[colName]=true
			columnsIndex[colName]=index
		}
	}
	var filteredRecords []patientRecord 
	shouldExit:=false
	for key, value:= range columnFound{
		if !value {
			fmt.Println("No such column found:",key)
			shouldExit=true
		}
	}
	if shouldExit {
		os.Exit(-1)
	}
	records,err:=csvReader.ReadAll()
	if ok!=nil {
		fmt.Println(err)
	}
	for _ , record:= range records{
		var (
			name=record[columnsIndex[nameColumn]]
			age,_ =strconv.ParseInt(record[columnsIndex[ageColumn]],10,32)
			disease=record[columnsIndex[diseaseColumn]]
			city=record[columnsIndex[diseaseColumn]]
		)
		if age>40 && strings.Contains(strings.ToLower(disease),"cancer") {
			var newRecord patientRecord= patientRecord{name:name , age:int(age), disease: disease, city: city }
			_=append(filteredRecords,newRecord)
		}
		
	}
	fmt.Println(len(filteredRecords))
}
	
	

	
	

