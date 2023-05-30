# SteelSeriesTest

To build : 

1. Install golang ( https://go.dev/doc/install ) event thought i'm pretty sure you already have it.
2. go run Server/SteelSeriesTest.go will launch the server.
3. go run Client/Client.go will launch the client and send a default http request to the server ( you will have to modify Client.go if you want to change the request )

You can also use the tests prepared in response_test.go. However make sure the server is not already running when you begin the tests.