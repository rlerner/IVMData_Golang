package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func basicAuth(username, password string) string { //function to base64 the user:pass string
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	userPtr := flag.String("user", "na", "HTTP Basic Auth USER") //Declare variables
	passPtr := flag.String("pass", "na", "HTTP Basic Auth PASS") //Declare variables
	hostPtr := flag.String("host", "na", "HOST+URL to request") //Declare variables
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

	fmt.Println("Let's Begin!")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //Ignore TLS configuration, as most IVM installs use a Self-signed cert

	client := &http.Client{} //Create the client to pass basicauth headers
	req, err := http.NewRequest("GET", host, nil) //Get Request
	req.Header.Set("Authorization", "Basic "+basicAuth(user, pass)) //set headers for basicauth using function from before
	resp, err := client.Do(req) //execute GET

	if err != nil {
		fmt.Println(err)
	} else {
		if _, err := os.Stat("sites.json"); os.IsNotExist(err) {
			file, err := os.Create("sites.json")
			if err != nil {
				fmt.Printf("Cannot create file for sites.json %s\n", err)
			} else {
				io.Copy(file, resp.Body)
				defer file.Close()
			}
		} else {
			var answer string
			fmt.Print("This file currently exists, do you wish to overwrite (y/N)?")
			fmt.Scanln(&answer)
			if answer == "y" || answer == "Y" {
				file, err := os.Create("sites.json")
				if err != nil {
					fmt.Printf("Cannot create file for sites.json %s\n", err)
				} else {
					io.Copy(file, resp.Body)
					defer file.Close()
				}
			} else {
				fmt.Println("You've stopped the program from replacing sites.json, aborting")
			}

		}

	}
	defer resp.Body.Close()
}
