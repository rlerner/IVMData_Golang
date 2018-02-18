package main

import (
	"fmt"
	"net/http"
	"encoding/base64"
	"flag"
	"os"
	"io"
	"strings"
)

func jsonToDisk(hostName string, userName string, password string, fileOutput string) {
	client := &http.Client{}
	req, requestError := http.NewRequest("GET",hostName,nil)
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(userName+":"+password)))
	resp, responseError := client.Do(req)
	
	if requestError != nil {
		fmt.Println("Cannot GET " + hostName + " %s\n",requestError)
	} else {
		if (responseError != nil) {
			fmt.Println("Error with response %s\n",responseError)
		}

		file, fileError := os.Create(fileOutput)
		if fileError != nil {
			fmt.Printf("Cannot create file for "+fileOutput+" %s\n", fileError)
		} else {
			io.Copy(file, resp.Body)
			defer file.Close()
		}
	}
	defer resp.Body.Close()
}

func main() {
	userPtr := flag.String("user","na","HTTP Basic Auth USER")
	passPtr := flag.String("pass","na","HTTP Basic Auth PASS")
	hostPtr := flag.String("host","na","HOST+URL to request")
	flag.Parse()

	var input string
	var user string = *userPtr
	if *userPtr == "na" {
		fmt.Print("A user is required: ")
		fmt.Scanln(&input)
		user = input
	}
	var pass string = *passPtr
	if *passPtr == "na" {
		fmt.Print("A password is required: ")
		fmt.Scanln(&input)
		pass = input
	}
	var host string = *hostPtr
	if *hostPtr == "na" {
		fmt.Print("A hostname+URL is required: ")
		fmt.Scanln(&input)
		host = input
	}


	var endPoints[8]string
	endPoints[0] = "asset_groups"
	endPoints[1] = "users"
	endPoints[2] = "tags"
	endPoints[3] = "administration/info"
	endPoints[4] = "administration/license"
	endPoints[5] = "scan_engines"
	endPoints[6] = "report_templates"
	endPoints[7] = "discovery_connections"

	for i := 0; i < len(endPoints); i++ {
		jsonToDisk(host+endPoints[i],user,pass,strings.Replace(endPoints[i],"/","-",-1)+".json")
	}
	
	fmt.Println("end")
}
