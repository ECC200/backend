package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var App *firebase.App
var StorageClient *storage.Client

func Initialize() {
	// サービスアカウントキーのファイルパスを指定
	opt := option.WithCredentialsFile("")
	var err error
	App, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Firebase Storageクライアントを初期化
	ctx := context.Background()
	StorageClient, err = storage.NewClient(ctx, option.WithCredentialsFile(""))
	if err != nil {
		log.Fatalf("error initializing storage client: %v\n", err)
	}
}
