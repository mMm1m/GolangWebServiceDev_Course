package main

import "testing"

func Test(t *testing.T) {
	server := TestingServer(ACCESS_TOKEN)
	defer server.close()

}
