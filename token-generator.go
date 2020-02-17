package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jws"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var cs *jws.ClaimSet
var hdr *jws.Header
var privateKey *rsa.PrivateKey

const authUrl = "https://www.googleapis.com/oauth2/v4/token"

func main() {

	port := flag.String("port", "8080", "port of the server")
	jsonFile := flag.String("file", "", "full path of json file")
	flag.Parse()

	file := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if *jsonFile != "" {
		file = *jsonFile
	}

	if file == "" {
		log.Fatal("Provide a json file either in the GOOGLE_APPLICATION_CREDENTIALS env variable or with the -file param ")
	}

	b, _ := ioutil.ReadFile(file)

	fromJSON, err := google.JWTConfigFromJSON(b)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Impossible to read the json file %s as a service account secret file", file)
	}

	cs = &jws.ClaimSet{
		Iss: fromJSON.Email,
		Sub: fromJSON.Email,
		Aud: authUrl,
	}
	hdr = &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     fromJSON.PrivateKeyID,
	}

	privateKey = ParseKey(fromJSON.PrivateKey)

	http.HandleFunc("/", tokenGenerator)
	http.ListenAndServe(":"+*port, nil)

}

func ParseKey(key []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil
	}
	return parsed
}

func tokenGenerator(w http.ResponseWriter, r *http.Request) {
	iat := time.Now()
	exp := iat.Add(time.Hour)

	aud := r.URL.Query().Get("aud")

	if aud == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "audience (host name of the service to reach) must be specified. add &aud=<HOST> to your request")
	}

	cs.PrivateClaims = map[string]interface{}{"target_audience": aud}
	cs.Iat = iat.Unix()
	cs.Exp = exp.Unix()

	msg, err := jws.Encode(hdr, cs, privateKey)
	if err != nil {
		fmt.Println(fmt.Errorf("google: could not encode JWT: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}

	f := url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {msg},
	}

	res, err := http.PostForm(authUrl, f)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}
	c, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	}

	type resIdToken struct {
		IdToken string `json:"id_token"`
	}

	id := &resIdToken{}

	json.Unmarshal(c, id)

	fmt.Fprintf(w, id.IdToken)
}
