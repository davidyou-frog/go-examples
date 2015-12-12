This is the program source that gets the token for API from the command line.

### Prepare the workspace

Set the GOPATH environment variable to your working directory.
Get the Gmail API Go client library and OAuth2 package using the following commands:

    $ go get -u google.golang.org/api/gmail/v1
    $ go get -u golang.org/x/oauth2/...