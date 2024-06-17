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
	opt := option.WithCredentialsFile("C:\\Users\\2220414\\Desktop\\システム開発\\care-connect-eba8d-firebase-adminsdk-cz4pl-eb03b9ff28.json")
	var err error
	App, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}
