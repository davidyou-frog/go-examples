package main

import (
    "os" 
    "fmt"
    "log"
	"io/ioutil"
	"encoding/json"
)

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"
)

type GmailTokenCreator struct {

	config *oauth2.Config
	code   string
	token  *oauth2.Token
	
}

func NewGmailTokenCreator( out_file string) GmailTokenCreator {

	return GmailTokenCreator{}
}

func (t *GmailTokenCreator) GetConfig( client_secret_file string ) *oauth2.Config {

	b, err := ioutil.ReadFile( client_secret_file )
    if err != nil {
        log.Fatalf("Unable to read client_secret.json: \n    %v", err)
    }

	t.config, err = google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
    if err != nil {
        log.Fatalf("Unable to parse client_secret.json: \n    %v", err)
    }
	
	return t.config
	
}

func (t *GmailTokenCreator) GetAuthCodeURL() string {

    if t.config == nil {
	    log.Fatalf("config.GmailTokenCreator is nil!\n")
	}
	
	return t.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	
}

func (t *GmailTokenCreator) SetAuthCode( code string ) {

    t.code = code
	
    tok, err := t.config.Exchange(oauth2.NoContext, t.code )
    if err != nil {
        log.Fatalf("Unable to retrieve token from web %v", err)
    }
	
	t.token = tok
}


func (t *GmailTokenCreator) SaveToken( token_file string ) *oauth2.Token {

    f, err := os.Create(token_file)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
     }
     defer f.Close()
	 
     json.NewEncoder(f).Encode(t.token)
	
    return t.token
}

const msg_intro = `get-token - create the token for google api
    from client_secret.json
    to   token.json`

const (
    in_file  = "client_secret.json"
    out_file = "token.json"
)	
	
func main() {

	fmt.Println( msg_intro )
	
	gtc := NewGmailTokenCreator( "token.json" )
	
	gtc.GetConfig( in_file )
	url := gtc.GetAuthCodeURL()
	
	fmt.Printf( "    Go to the following link in your browser" )
	fmt.Printf( "\n\n%s\n\n", url )

    fmt.Printf( "    type the authorization code: " )
	var code string 
    if _, err := fmt.Scan(&code); err != nil {
        log.Fatalf("Unable to read authorization code %v", err)
    }
	
	gtc.SetAuthCode( code )
	
	fmt.Printf("    Saving credential file to: %s\n", out_file )
	gtc.SaveToken( out_file )
	
}

