package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(os.Getenv("MY_JWT_TOKEN_KEY"))

func homePage(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Request to server
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:9000/", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}

	fmt.Fprintf(w, string(body))
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Dany M Pradana"
	// Set expired 1 minute
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Printf("Something error: %s", err)
		return "", err
	}

	return tokenString, nil
}

func handleRequest() {
	http.HandleFunc("/", homePage)

	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	// Check jika token key is not assigned
	if len(mySigningKey) == 0 {
		fmt.Println("JWT token key is missing")
		os.Exit(1)
	}
	// tokenString, err := GenerateJWT()
	// if err != nil {
	// 	fmt.Println("Error generating token string: ", err.Error())
	// }

	// fmt.Println("My token string: ", tokenString)

	// REST API
	handleRequest()
}
