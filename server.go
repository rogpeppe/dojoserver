package main

import (
	"fmt"
	//"net"
	"bytes"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// etcd endpoint
const endpoint = "http://localhost:2379/v2/keys/"

func main() {
	// Define routes and handlers
	r := mux.NewRouter()
	r.HandleFunc("/addresses", GetAddresses).Methods("GET")
	r.HandleFunc("/addresses/new", PutAddress).Methods("PUT")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

// GetAddresses makes a GET request to etcd to list all addressed that have been
// sent and will ideally return a key / value list of them for requests made to
// this endpoint. The raw response is printed out to the writer.
func GetAddresses(w http.ResponseWriter, r *http.Request) {
	addresses := getEtcd()

	// All this does is print the response from etcd at the moment we should
	// return probably a key / value list of the addresses
	fmt.Fprint(w, addresses)
}

// PutAddress make a PUT request to etcd to store the address against a teamname
// and prints out the raw response to the writer.
func PutAddress(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	teamname := r.Form.Get("teamname")

	// We don't have the address at this point, we could get it from the request
	// object I think. Currently the type to be passed through is a string.
	address := "needs to be set here!!!"

	resp := putEtcd(teamname, address)

	// Same as the get, this will still just print the JSON generated by etcd
	fmt.Fprint(w, resp)
}

// getEctd makes the request to etcd to return a raw JSON response of all the
// keys. I haven't done anything with the error handling as of yet.
func getEtcd() string {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

// putEtcd makes the internal request to etcd to PUT an address as a value to a
// teamname which is passed through as a PUT param.
func putEtcd(teamname, address string) string {
	requestUrl := generateEndpoint(teamname)

	form := url.Values{}
	form.Add("value", address)

	req, err := http.NewRequest("PUT", requestUrl, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

// generateEndpoint simply concats the endpoint const with the teamname passed
// through as a PUT param and returns the buffer as a string.
func generateEndpoint(teamname string) string {
	var b bytes.Buffer
	b.WriteString(endpoint)
	b.WriteString(teamname)

	return b.String()
}
