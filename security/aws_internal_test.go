package security

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/urfave/cli"
	"github.com/stretchr/testify/assert"
)

func TestAppliesIfAwsInternalFlagSet(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	serverAddress := "/security/a-role"
	provider := AwsInternalProvider{ServerAddress: "http://localhost:8888/security/"}
	go StartServer(serverAddress)
	err := provider.Apply(awsDefaultContext())

	assert.Nil(t, err)
	assert.Equal(t, os.Getenv("AWS_ACCESS_KEY_ID"), "fakeaccesskeyid")
	assert.Equal(t, os.Getenv("AWS_SECRET_ACCESS_KEY"), "fakesecretaccesskey")
}

func StartServer(serverAddress string) {
	http.HandleFunc(serverAddress, respondWithJson) // set router
	err := http.ListenAndServe(":8888", nil)          // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func respondWithJson(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Code":            "Success",
		"LastUpdated":     "2015-11-04T02:53:03Z",
		"Type":            "AWS-HMAC",
		"AccessKeyId":     "fakeaccesskeyid",
		"SecretAccessKey": "fakesecretaccesskey",
		"Token":           "AQoDXXdzECQa8AN8q3pJ94WoT9PACggVSrVQP/YEvLWuHC09mTZTtd8ruywhFB/6Z//+b9xyQSzxZxRspz+TymnIct3D8Y2hg3KqAic/JtwANId9Z5gRiAIlaKJsNWqNMIOF6MM1ac7PN48xRICFJ1cc+w4lVuOhK0bmpkCu8n0+Gx4kClZBEqdZHqXMoTzYNnxg4x1rk1AgOYyir0u9rxRfoxjZ0h/kE68P6OQ9ra/aIUANYJMdbBZ/6KWXEiX3Elffu9YE3VGIAAvp+itFdCFy64v2agMZxvazoiKfnKmV44IwAXfOituZB80/eEU+/fNZYLEEteM0V8y5Qt0ipw7zjFS/IJ3YceR361q/5Mut8WAEXfprIO/5BrIHKcjUXuUm658Toi9dcPXUNhDdYmRLoLYU81PiNPYatJRhxTpekKhUaMv5r1Ikm9tixqy5TwAt8G0P3nJ3rP8Z4ZympRKqS0nK41ni6XTYy0x1z/ss3BkUd8Z27isBy87Ngrp/7CxYnnZ+ynhx3MvaZyxdw8a9Jbticg9X6mYVsBQPiColquVC92Ei8K6jpf+6lRTEudL8U3+UKEerPyCpULdwP5KeOsB2MFzXtedB4IJs1qrDrxzIG9b1AADtfb+GcZ+bxLgODvPLdQqCf1GQ8EeKGP6SiZIT11/EOCjZIO7l5bEF",
		"Expiration":      "2015-11-04T09:04:47Z",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func awsDefaultContext() *cli.Context {
	app := cli.NewApp()
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.String("security", "aws-internal", "apply security")
	set.String("aws-role", "a-role", "aws role")
	context := cli.NewContext(app, set, nil)
	return context
}
