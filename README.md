# analyticsdumper
Go code that queries GA data via the Google Core Reporting API using two-legged service account Google OAuth2 authentication

##### overview
demo code that shows:
* service account authentication with a Google API in Go
* how to pull pageview count from last 24 hours via Core Reporting API in Go

##### HOW-TO
* create a new app in the Google Cloud Console
 * see http://code.rickt.org/2014/03/how-to-download-google-analytics-data.html
  * get your app's private key from the Google Cloud Console
  * get your app's client secrets JSON file from the Google Cloud Console
* install required Go libraries
 * `go get golang.org/x/oauth2`
 * `go get golang.org/x/oauth2/jwt`
 * `go get google.golang.org/api/analytics/v3`
* download `analyticsdumper.go`
* edit `analyticsdumper.go`
 * change `ga*` vars in code as noted in comments to values from your client secrets JSON file
* `go build analyticsdumper.go && ./analyticsdumper`
* profit!
* example output:
 $ go build analyticsdumper.go && ./analyticsdumper
 [1399856] pageviews for www.redacted.com (UA-redacted-1) from 2016-04-06 to 2016-04-07.
