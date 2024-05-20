package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// User structure for Firestore data
// Firestoreのデータをマッピングするための構造体
type User struct {
	ID   string `json:"id"`   // ドキュメントID
	Name string `json:"name"` // ユーザーの名前
	Age  int    `json:"age"`  // ユーザーの年齢
}

func main() {
	// Use the service account file
	// サービスアカウントのJSONファイルを使ってFirebaseアプリを初期化
	sa := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	// Initialize Firestore client
	// Firestoreクライアントを初期化
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
	}
	// 関数終了時にクライアントを閉じる
	defer client.Close()

	// Ginのデフォルトのルーターを作成
	router := gin.Default()

	// Endpoint to get users from Firestore
	// Firestoreからユーザー情報を取得するエンドポイント
	router.GET("/users", func(c *gin.Context) {
		var users []User
		// Firestoreの`users`コレクションからドキュメントを取得
		iter := client.Collection("users").Documents(context.Background())
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				// ドキュメントの取得が完了した場合にループを終了
				break
			}
			if err != nil {
				log.Fatalf("error iterating documents: %v", err)
			}

			// ドキュメントのデータをUser構造体にマッピング
			user := User{
				ID:   doc.Ref.ID,                     // ドキュメントID
				Name: doc.Data()["name"].(string),    // 名前フィールド
				Age:  int(doc.Data()["age"].(int64)), // 年齢フィールド
			}
			// ユーザーリストに追加
			users = append(users, user)
		}
		// ユーザーリストをJSON形式でレスポンスとして返す
		c.JSON(http.StatusOK, users)
	})
	// HTTPサーバーをポート8080で起動(修正が必要)
	router.Run(":8080")
}
