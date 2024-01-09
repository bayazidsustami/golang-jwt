package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenParse(t *testing.T) {
	authorizationHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJNeSBTaW1wbGUgSldUIEFQUCIsImV4cCI6MTcwNDc3OTE2NCwidXNlcm5hbWUiOiJiYXlhemlkIiwiZW1haWwiOiJ0ZXJwYWxtdXJhaEBnbWFpbC5jb20iLCJHcm91cCI6ImFkbWluIn0.d-LzhJ_Wt0G3r_r0lO0LNlw2p7yyOTYbPnP-NpR67cc"
	token := strings.Replace(authorizationHeader, "Bearer ", "", -1)

	fmt.Println(token)
}
