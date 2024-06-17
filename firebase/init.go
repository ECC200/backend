package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var App *firebase.App

func Initialize() {
	// サービスアカウントキーのファイルパスを指定(githubに挙げるときは記述しない)
	opt := option.WithCredentialsFile("../../care-connect-eba8d-firebase-adminsdk-cz4pl-579f118464.json")
	var err error
	App, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}
