package main


import (
    "os" 
    "fmt"
     "log"
 	"io/ioutil"
 	"encoding/json"
	"encoding/base64"
//    "net/http"
)

import (
    "golang.org/x/net/context" 
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"
	
	"github.com/jaytaylor/html2text"
)

var DefaultAcceptAttachmentMimeTypes = []string{
    "application/octet-stream",
	"image/png",
	"image/jpg",
	"image/jpeg",
	"image/gif",
}

type GmailManager struct {

    email      string
    ctx        context.Context

	config     *oauth2.Config
	token      *oauth2.Token
	
	srv        *gmail.Service
	msgs       []*GmailMessage 
	
	acceptAttachmentMimeTypes []string
}

func NewGmailManager( email string ) GmailManager {

    gm     := GmailManager{} 
	
    gm.ctx   = context.Background()
	gm.email = email
	gm.msgs  = []*GmailMessage{} 
	gm.acceptAttachmentMimeTypes = DefaultAcceptAttachmentMimeTypes
	
	return gm
}

func (gm *GmailManager) GetConfig( client_secret_file string ) *oauth2.Config {

	b, err := ioutil.ReadFile( client_secret_file )
    if err != nil {
        log.Fatalf("Unable to read %s: \n    %v", client_secret_file, err)
    }

	gm.config, err = google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
    if err != nil {
        log.Fatalf("Unable to parse %s: \n    %v", client_secret_file, err)
    }
	
	return gm.config

}

func (gm *GmailManager) LoadToken( token_file string ) *oauth2.Token {

    f, err := os.Open(token_file)
    if err != nil {
	    log.Fatalf("Unable to read %s: \n    %v", token_file, err)
    }
	
    gm.token = &oauth2.Token{}
    err = json.NewDecoder(f).Decode(gm.token)
    defer f.Close()
    if err != nil {
        log.Fatalf("Unable to parse %s: \n    %v", token_file, err)
    }
	
    return gm.token
}

func (gm *GmailManager) GetService() *gmail.Service {

    client := gm.config.Client(gm.ctx, gm.token ) 
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}
	
	gm.srv = srv
	return gm.srv
}

func (gm *GmailManager) BuildService( client_secret_file, token_file string ) {

    gm.GetConfig( client_secret_file ) 
	gm.LoadToken( token_file ) 
	gm.GetService()
	
}

func (gm *GmailManager) GetMailList() []*GmailMessage {

    // req := srv.Users.Messages.List(gm.email).LabelIds( "INBOX", "UNREAD" ).Q( "subject:[build]" )
    req := gm.srv.Users.Messages.List(gm.email).LabelIds( "INBOX" )
    res, err := req.Do()
    if err != nil {
	    log.Fatalf("Unable to retrieve messages: %v", err)
    }

	for _, m := range res.Messages {

	    msg :=  NewGmailMessage( gm, m.Id ) 
	    gm.msgs = append( gm.msgs, &msg )

    }
	
	return gm.msgs
	
}

type GmailMessage struct {

	gm          *GmailManager
	
    Id             string
	Subject        string
	Sender         string
	Body           string
	
	Attachments []*GmailAttachment

}

func NewGmailMessage( gm *GmailManager, id string ) GmailMessage {

    m   :=  GmailMessage{}
    m.gm          =  gm	
	m.Id          =  id
	m.Attachments = []*GmailAttachment{}
	
	return m
}

func ( m *GmailMessage ) getHeaderValue( headers []*gmail.MessagePartHeader, name string) string {
	for _, one := range headers {
		if one.Name == name {
			return one.Value
		}
	}

	return ""
}

func ( m *GmailMessage ) getSubject( headers []*gmail.MessagePartHeader) string {
	return m.getHeaderValue( headers, "Subject")
}

func ( m *GmailMessage ) getSender( headers []*gmail.MessagePartHeader) string {
	return m.getHeaderValue( headers, "From")
}

func ( m *GmailMessage ) getBody(parts []*gmail.MessagePart) string {

	for _, part := range parts {
		if len(part.Parts) > 0 {
			return m.getBody(part.Parts)
		} else {
			if part.MimeType == "text/html" {
				return part.Body.Data
			}
		}
	}

	return ""
}

func ( m *GmailMessage ) isAcceptAttachmentMimeType(mime string) bool {

	for _, mt := range m.gm.acceptAttachmentMimeTypes {
		if mt == mime {
			return true
		}
	}

	return false
}


func ( m *GmailMessage ) getAttachments(parts []*gmail.MessagePart) {

	for _, part := range parts {
		if len(part.Parts) == 0 {
		    
			fmt.Printf( "part.MimeType = %s\n", part.MimeType )

			if m.isAcceptAttachmentMimeType( part.MimeType ) {

	            a  :=  NewGmailAttachment( m, 
				                           part.Body.AttachmentId,
                                           part.MimeType,
					                       part.Filename )
										   
	            m.Attachments = append( m.Attachments, &a )
            } 
		}
	}
	
}

func ( m *GmailMessage ) GetMail() {

    req := m.gm.srv.Users.Messages.Get( m.gm.email, m.Id )
	res, err := req.Format("full").Do()
    if err != nil {
	    log.Fatalf("Unable to get messages: %v", err)
    }
	
	m.Subject = m.getSubject( res.Payload.Headers )
	m.Sender  = m.getSender ( res.Payload.Headers )
	m.Body    = m.getBody   ( res.Payload.Parts   )
	
	m.getAttachments ( res.Payload.Parts   ) 
	
}


func ( m *GmailMessage) GetBodyHTML() string {

    data,_ :=  base64.URLEncoding.DecodeString(m.Body)
//	html   := base64.StdEncoding.EncodeToString(data)

	return string(data)
}

func ( m *GmailMessage) GetBodyTEXT() string {

    text, _ := html2text.FromString( m.GetBodyHTML() )
	
	return string(text)
	
}

type GmailAttachment struct {

	msg        *GmailMessage
	
    Id          string
	MimeType    string
	Filename    string
	
	Data        []byte
	
}

func NewGmailAttachment( m *GmailMessage, id, mimetype, filename string ) GmailAttachment {

    a   := GmailAttachment{} 
	a.msg      = m 
	a.Id       = id
	a.MimeType = mimetype
	a.Filename = filename 
	
	return a
	
}

func ( a *GmailAttachment) GetAttachment() {

    m   := a.msg 
    gm  := m.gm
    req := gm.srv.Users.Messages.Attachments.Get( gm.email, m.Id, a.Id )
	res, err := req.Do()
    if err != nil {
	    log.Fatalf("Unable to retrieve attachment: %v", err)
    }
	
	a.Data, _ = base64.URLEncoding.DecodeString( res.Data )
	
}

func ( a *GmailAttachment) SaveAs( filename string ) {

    f, err := os.Create(filename)
    if err != nil {
        log.Fatalf("Unable to save as : %v", err)
     }
     defer f.Close()
	 
	 f.Write( a.Data )
	 
}

func ( a *GmailAttachment) Save() {

    a.SaveAs( a.Filename )

}

