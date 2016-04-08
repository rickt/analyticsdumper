package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/analytics/v3"
	"io/ioutil"
	"log"
	"time"
)

// constants
const (
	datelayout string = "2006-01-02" // date format that Core Reporting API requires
)

// globals that you DON'T need to change
var (
	enddate   string = time.Now().Format(datelayout)                          // set end query date to today
	startdate string = time.Now().Add(time.Hour * 24 * -1).Format(datelayout) // set start query date to yesterday
	metric    string = "ga:pageviews"                                         // GA metric that we want
	scope     string = "https://www.googleapis.com/auth/analytics.readonly"   // GA API scope
	tokenurl  string = "https://accounts.google.com/o/oauth2/token"           // (json:"token_uri") Google oauth2 Token URL
)

// globals that you DO need to change

// populate these with values from the JSON secretsfile obtained from the Google Cloud Console specific to your app)
// example secretsfile JSON:
// {
// 	"web": {
// 		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
// 		"token_uri": "https://accounts.google.com/o/oauth2/token",
// 		"client_email": "blahblahblahblah@developer.gserviceaccount.com",
// 		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/blahblahblahblah@developer.gserviceaccount.com",
// 		"client_id": "blahblahblahblah.apps.googleusercontent.com",
// 		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs"
// 	}
// }
var (
	// CHANGE THESE!!!
	gaServiceAcctEmail  string = "blahblahblahblah@developer.gserviceaccount.com" // (json:"client_email") email address of registered application
	gaServiceAcctPEMKey string = "/tmp/analyticsdumper.pem"                       // full path to private key file (PEM format) of your application from Google Cloud Console
	gaTableID           string = "ga:NNNNNNNN"                                    // namespaced profile (table) ID of your analytics account/property/profile
)

// func: main()
// the main function.
func main() {
	// load up the registered applications private key
	pk, err := ioutil.ReadFile(gaServiceAcctPEMKey)
	if err != nil {
		log.Fatal("Error reading GA Service Account PEM key -", err)
	}
	// create a jwt.Config that we will subsequently use for our authenticated client/transport
	// relevant docs for all the oauth2 & json web token stuff at https://godoc.org/golang.org/x/oauth2 & https://godoc.org/golang.org/x/oauth2/jwt
	jwtc := jwt.Config{
		Email:      gaServiceAcctEmail,
		PrivateKey: pk,
		Scopes:     []string{scope},
		TokenURL:   tokenurl,
	}
	// create our authenticated http client using the jwt.Config we just created
	clt := jwtc.Client(oauth2.NoContext)
	// create a new analytics service by passing in the authenticated http client
	as, err := analytics.New(clt)
	if err != nil {
		log.Fatal("Error creating Analytics Service at analytics.New() -", err)
	}
	// create a new analytics data service by passing in the analytics service we just created
	// relevant docs for all the analytics stuff at https://godoc.org/google.golang.org/api/analytics/v3
	ads := analytics.NewDataGaService(as)
	// w00t! now we're talking to the core reporting API. the hard stuff is over, lets setup a simple query.
	// setup the query, call the Analytics API via our analytics data service's Get func with the table ID, dates & metric variables
	gasetup := ads.Get(gaTableID, startdate, enddate, metric)
	// send the query to the API, get a big fat gaData back.
	gadata, err := gasetup.Do()
	if err != nil {
		log.Fatal("API error at gasetup.Do() -", err)
	}
	// print out some nice things
	fmt.Printf("%s pageviews for %s (%s) from %s to %s.\n", gadata.Rows[0], gadata.ProfileInfo.ProfileName, gadata.ProfileInfo.WebPropertyId, startdate, enddate)
	return
}
