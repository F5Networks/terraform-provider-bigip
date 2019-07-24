package bigip

import (
        "fmt"
        "log"
        "net/http"
        "io/ioutil"
        "crypto/tls"
        "github.com/f5devcentral/go-bigip"
        "strings"
        "github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipAs3() *schema.Resource {
        return &schema.Resource{
                Create: resourceBigipAs3Create,
                Read:   resourceBigipAs3Read,
                Update: resourceBigipAs3Update,
                Delete: resourceBigipAs3Delete,
                Exists: resourceBigipAs3Exists,
                Importer: &schema.ResourceImporter{
                        State: schema.ImportStatePassthrough,
                },

                Schema: map[string]*schema.Schema{

                        "as3_json": {
                                Type:        schema.TypeString,
                                Required:    true,
                                Description: "AS3 json",
                        },
                        "tenant_name": {
                                Type:        schema.TypeString,
                                Optional:    true,
                                Description: "unique identifier for resource",
                        },
                },
        }
}

func resourceBigipAs3Create(d *schema.ResourceData, meta interface{}) error {
        client_bigip := meta.(*bigip.BigIP)

        as3_json := d.Get("as3_json").(string)
        name := d.Get("tenant_name").(string)
        log.Printf("[INFO] as3_json is :%s", as3_json )
        tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
        client := &http.Client{Transport: tr}
        url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
        req, err := http.NewRequest("POST", url, strings.NewReader(as3_json))
        if err != nil {
	   return fmt.Errorf("Error while creating http request with AS3 json:%v",err)
        }
        req.SetBasicAuth(client_bigip.User,client_bigip.Password)
        req.Header.Set("Accept", "application/json")
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
	log.Printf("[INFO] as3 resp in update call is :%+v", resp)
        body, err := ioutil.ReadAll(resp.Body)
        bodyString := string(body)
        if ( resp.Status != "200 OK" ||  err != nil)  {
           defer resp.Body.Close() 
	   return fmt.Errorf("Error while Sending/Posting http request with AS3 json :%s  %v",bodyString,err )
        } 
        log.Printf("[INFO] as3 resp in create call is :%+v", resp)
        log.Printf("[INFO] as3 err in create call is :%v", err )

        defer resp.Body.Close()
        d.SetId(name)
        return resourceBigipAs3Read(d,meta)
}
func resourceBigipAs3Read(d *schema.ResourceData,meta interface{}) error {
        client_bigip := meta.(*bigip.BigIP)
        as3_json := d.Get("as3_json").(string)
        name := d.Id()
        log.Printf("[INFO] as3_json in read is :%s", as3_json )
        log.Println("[INFO] tenant name is :%s",name)

        tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
        client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
        req, err := http.NewRequest("GET", url, strings.NewReader(as3_json))
        if err != nil {
           return fmt.Errorf("Error while creating http request with AS3 json:%v",err)
        }
        req.SetBasicAuth(client_bigip.User,client_bigip.Password)
        req.Header.Set("Accept", "application/json")
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        body, err := ioutil.ReadAll(resp.Body)
        bodyString := string(body)
        log.Printf("[INFO] as3 resp in read call is :%+v", resp)
        log.Printf("[INFO] as3 bodystring in read is :%+v",bodyString)
        if ( resp.Status != "200 OK" ||  err != nil)  {
           defer resp.Body.Close()
           return fmt.Errorf("Error while Sending/fetching http request with AS3 json :%s  %v",bodyString,err )
        }

        log.Printf("[INFO] as3 resp in read call is :%v", resp)
        log.Printf("[INFO] as3 err in read call is :%v", err)
        defer resp.Body.Close()
        return nil
}

func resourceBigipAs3Exists(d *schema.ResourceData,meta interface{}) (bool, error) {
        client_bigip := meta.(*bigip.BigIP)
        log.Printf("[INFO] Checking if As3 config exists in bigip ")

        tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
        client := &http.Client{Transport: tr}
        url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
           log.Printf("[ERROR] Error while creating http request with AS3 json: %v", err)
           return false, err
        }
        req.SetBasicAuth(client_bigip.User,client_bigip.Password)
        req.Header.Set("Accept", "application/json")
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        body, err := ioutil.ReadAll(resp.Body)
        bodyString := string(body)
        log.Printf("[INFO] as3 resp in Exists call is :%+v", resp)
        log.Printf("[INFO] as3 bodystring in Exists call is :%+v",bodyString)
        if ( resp.Status == "204 No Content" ||  err != nil)  {
           log.Printf("[ERROR] Error while checking as3resource present in bigip :%s  %v",bodyString,err )
	   defer resp.Body.Close()
           return false, err
        }
        defer resp.Body.Close()
        return true, nil
}


func resourceBigipAs3Update(d *schema.ResourceData, meta interface{}) error {
        client_bigip := meta.(*bigip.BigIP)
        as3_json := d.Get("as3_json").(string)
        log.Printf("[INFO] as3_json in update is :%s", as3_json )
        log.Printf("[INFO] in update cleint_bigip is :%v", client_bigip)
        tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true} }
        client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
        req, err := http.NewRequest("PATCH", url, strings.NewReader(as3_json))
        if err != nil {
           return fmt.Errorf("Error while creating http request with AS3 json:%v",err)
        }
        req.SetBasicAuth(client_bigip.User,client_bigip.Password)
        req.Header.Set("Accept", "application/json")
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        body, err := ioutil.ReadAll(resp.Body)
        bodyString := string(body)
        if ( resp.Status != "200 OK" ||  err != nil)  {
           return fmt.Errorf("Error while Sending/Posting http request with AS3 json :%s  %v",bodyString,err )
        }

        log.Printf("[INFO] as3 resp in update call is :%+v", resp)
        log.Printf("[INFO] as3 err in update call is :%v", err)
        defer resp.Body.Close()
        return resourceBigipAs3Read(d,meta)
}

func resourceBigipAs3Delete(d *schema.ResourceData, meta interface{}) error {
        client_bigip := meta.(*bigip.BigIP)
        as3_json := d.Get("as3_json").(string)
        log.Printf("[INFO] as3_json in update is :%s", as3_json )

        tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true} }
        client := &http.Client{Transport: tr}
	url := client_bigip.Host + "/mgmt/shared/appsvcs/declare"
        req, err := http.NewRequest("DELETE", url, strings.NewReader(as3_json))

        if err != nil {
           return fmt.Errorf("Error while creating http request with AS3 json:%v",err)
        }
        req.SetBasicAuth(client_bigip.User,client_bigip.Password)
        req.Header.Set("Accept", "application/json")
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        body, err := ioutil.ReadAll(resp.Body)
        bodyString := string(body)
        if ( resp.Status != "200 OK" ||  err != nil)  {
           return fmt.Errorf("Error while Sending/deleting http request with AS3 json :%s  %v",bodyString,err )
        }

        log.Printf("[INFO] as3 resp in delete call is :%v", resp)
        log.Printf("[INFO] as3 err in delete call is :%v", err)
        defer resp.Body.Close()
        d.SetId("")
        return nil
}

