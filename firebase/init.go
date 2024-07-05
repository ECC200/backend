package firebase

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var App *firebase.App
var StorageClient *storage.Client

func Initialize() {
	// サービスアカウントキーのファイルパスを指定
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	var err error
	App, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Firebase Storageクライアントを初期化
	ctx := context.Background()
	StorageClient, err = storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("error initializing storage client: %v\n", err)
	}
}

// 署名付きURLを生成
func GenerateSignedURL(bucket, object string) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: "your-service-account-email",
		PrivateKey:     []byte("your-private-key"),
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	}

	url, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return "", err
	}
	return url, nil
}
