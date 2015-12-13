This is the program source that gets the token for the Gmail API from the command line.

### Prepare the workspace

Set the GOPATH environment variable to your working directory.

Get the Gmail API Go client library and OAuth2 package using the following commands:

    $ go get -u google.golang.org/api/gmail/v1
    $ go get -u golang.org/x/oauth2/...
	
### Build

    $ cd $GOPATH/src/oauth/google/get_token/
	$ go build
	
### Run

get_token을 아래와 같이 실행한다. 

Run get_token as follow
	
	$ cd $GOPATH/src/oauth/google/get_token/
	$ ./get_token
	get-token - create the token for google api
    from client_secret.json
    to   token.json
    Go to the following link in your browser

    https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=414979728339-e3pjqvpm9852h96r1b9n5uhbrbl3vpua.apps.googleusercontent.com&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fgmail.readonly&state=state-token

    type the authorization code: 

If the get-token is excuted nomally, then you can see the URL and wait to type the authorization code.
