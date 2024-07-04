package firebase

import (
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var App *firebase.App
var StorageClient *storage.Client
var err error

func Initialize(c *gin.Context) (*firestore.Client, *storage.Client) {
	opt := option.WithCredentialsFile("/Users/cyrusman/Desktop/ECCコンピューター専門学校/Year-3/Y3-Sem1/ITプロジェクト開発/Team2024/care-connect-eba8d-firebase-adminsdk-cz4pl-066e5ec155.json")
	// サービスアカウントキーのファイルパスを指定
	App, err = firebase.NewApp(c, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Firestoreクライアントを初期化
	Client, err := App.Firestore(c)
	if err != nil {
		log.Fatalf("error :Failed to initialize Firestore: %v", err)
	}

	// Firebase Storageクライアントを初期化
	StorageClient, err = storage.NewClient(c, opt)
	if err != nil {
		log.Fatalf("error initializing storage client: %v\n", err)
	}

	return Client, StorageClient
}
