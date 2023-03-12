package main

import "net/http"

type MyJWTTransport struct {
	token string
	transport http.RoundTripper
}


/** So we needed a round tripper function, but nstead of writing our own all 
	we want to do was add a header to the default Roundtrip function. 
	So here we added a transport object to the jwt struct then made it of type http.RoundTripper
	then in main.go we passed in the default transport that is already of type RoundTripper.
	That gave us access to the default RoundTrip methoud for the Roundtripper struct
	so in out function all we had to modify the in the request object, then run the default function.  
**/
func (m MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.token != ""{
		req.Header.Add("Authorization", "Bearer "+m.token)
	}
	return m.transport.RoundTrip(req)
}
