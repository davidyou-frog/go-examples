package main

import (
    "fmt"
)
 
const msg_intro = `get-gmail - receive a mail from gmail by the gmail api`

const (
    in_file  = "client_secret.json"
    out_file = "token.json"
)	
	
func main() {

	fmt.Println( msg_intro )

//    email := "service.davidyou@gamil.com"
    email := "me"	
	
	gm  := NewGmailManager( email )
	gm.BuildService       ( "client_secret.json" , "token.json" )

	list := gm.GetMailList()
	
	for _, m := range list {
	
		m.GetMail()
	    fmt.Printf( "ID = %s\n", m.Id )
	    fmt.Printf( "   Subject   [%s]\n", m.Subject )
		
//		fmt.Printf( "   Sender    [%s]\n", m.Sender  )
//		fmt.Printf( "   Body      [%s]\n", m.Body    )
//		fmt.Printf( "   Body(HTML)[%s]\n", m.GetBodyHTML() )
//		fmt.Printf( "   Body(TEXT)[%s]\n", m.GetBodyTEXT() )

    	for _, a := range m.Attachments {
    	
    	    fmt.Printf( "    Attachment Id = %s\n", a.Id )
    	    fmt.Printf( "    Attachment MimeType = %s\n", a.MimeType )
    	    fmt.Printf( "    Attachment Filename = %s\n", a.Filename )
    		
			a.GetAttachment()
			a.Save()
//			fmt.Printf( "    Attachment Data = %v\n", a.Data )
    	}

    }
	
}

