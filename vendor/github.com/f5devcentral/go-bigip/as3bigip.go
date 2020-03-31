package bigip

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
)

const as3SchemaLatestURL = "https://raw.githubusercontent.com/F5Networks/f5-appsvcs-extension/master/schema/latest/as3-schema.json"

type as3Validate struct {
	as3SchemaURL    string
	as3SchemaLatest string
}

func ValidateAS3Template(as3ExampleJson string) bool {
	myAs3 := &as3Validate{
		as3SchemaLatestURL,
		"",
	}
	err := myAs3.fetchAS3Schema()
	if err != nil {
		fmt.Errorf("As3 Schema Fetch failed: %s", err)
		return false
	}

	schemaLoader := gojsonschema.NewStringLoader(myAs3.as3SchemaLatest)
	//schemaLoader := gojsonschema.NewReferenceLoader("file:///Users/chinthalapalli/go/src/github.com/Practice/as3-schema-3.13.2-1-cis.json")
	documentLoader := gojsonschema.NewStringLoader(as3ExampleJson)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Errorf("%s", err)
		return false
	}
	if !result.Valid() {
		log.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
		return false
	}
	return true
}

func (as3 *as3Validate) fetchAS3Schema() error {
	res, resErr := http.Get(as3.as3SchemaURL)
	if resErr != nil {
		log.Printf("Error while fetching latest as3 schema : %v", resErr)
		return resErr
	}
	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Unable to read the as3 template from json response body : %v", err)
			return err
		}
		defer res.Body.Close()
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			log.Printf("Unable to unmarshal json response body : %v", err)
			return err
		}
		jsonMap["$id"] = as3SchemaLatestURL
		byteJSON, err := json.Marshal(jsonMap)
		if err != nil {
			log.Printf("Unable to marshal : %v", err)
			return err
		}
		as3.as3SchemaLatest = string(byteJSON)
		return err
	}
	return nil
}
