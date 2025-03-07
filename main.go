package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"iis_server/apiq"
	"iis_server/config"
	"iis_server/scheduler"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot setup program environment\n")
		os.Exit(1)
	}
	scheduler.Start()
}

func setup() error {
	_ = os.Mkdir(config.UPLOAD_FOLDER, 0755)
	_ = os.Mkdir(config.TMP_FOLDER, 0755)
	_ = os.Mkdir(config.XML_SCHEMAS_FOLDER, 0755)

	logger, err := zap.NewDevelopment()
	if err != nil {
		zap.ReplaceGlobals(zap.NewNop())
	} else {
		zap.ReplaceGlobals(logger)
		zap.S().Infof("Console logger setup")
	}

	if err := godotenv.Load(); err != nil {
		zap.S().DPanicf("failed to load .env err = %s", err.Error())
	}
	config.RapidApiKey = os.Getenv("RAPIDAPI_KEY")

	//	testApi()
	TestJsonToXml()
	return nil
}

func testApi() {
	api, err := apiq.IgApiFactory()
	if err != nil {
		panic(err)
	}

	value, err := api.GetUsernameByUserId("18527")
	fmt.Printf("err: %v\n", err)
	fmt.Printf("value: %v\n", value)
}

func TestJsonToXml() {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()

	// Read file content
	byteValue, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Decode JSON
	var data apiq.UserInfo
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Convert to XML
	xmlData, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error converting to XML:", err)
		return
	}

	// Print XML
	fmt.Println(string(xmlData))

	// Optionally, write to an XML file
	err = os.WriteFile("data.xml", xmlData, 0644)
	if err != nil {
		fmt.Println("Error writing XML file:", err)
	}
}
