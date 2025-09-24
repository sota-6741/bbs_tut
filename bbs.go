package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"
)

const logFile = "logs.json" // データの保存先

// Log 掲示板に保存するデータを構造体で定義
type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
	CTime int64  `json:"ctime"`
}

// ログファイルの書き込み
func saveLogs(logs []Log) {
	// JSONにエンコード
	bytes, error := json.Marshal(logs)
	if error != nil {
		fmt.Printf("ログデータのJSON変換に失敗しました: %v\n", error)
		return
	}
	// ファイルへ書き込む
	os.WriteFile(logFile, bytes, 0644)
}

// ファイルからログファイルの読み込み
func loadLogs() []Log {
	// ファイルを開く
	text, error := os.ReadFile(logFile)
	if error != nil {
		fmt.Printf("ファイルの読み込みに失敗しました: %v\n", error)
		return make([]Log, 0)
	}
	// JSONをパース
	var logs []Log
	json.Unmarshal([]byte(text), &logs)
	return logs
}

func showHandler(writer http.ResponseWriter, request *http.Request) {
	// ログを読み出してHTMLを生成
	var htmlLog string = ""
	var logs []Log = loadLogs() // データを読み出す
	for _, logEntry := range logs {
		htmlLog += fmt.Sprintf(
			"<p>(%d) <span>%s</span>: %s --- %s</p>",
			logEntry.ID,
			html.EscapeString(logEntry.Name),
			html.EscapeString(logEntry.Body),
			time.Unix(logEntry.CTime, 0).Format("2006/1/2 15:04"))
	}
	// HTML全体を出力
	htmlBody := "<html><head><style>" +
		"p {border: 1px solid silver; padding; 1em;}" +
		"span {background-color: #eef}" +
		"</style></head><body><h1>BBS</h1>" +
		getForm() + htmlLog + "</body></html>"
	writer.Write([]byte(htmlBody))
}

// フォームから送信された内容を書き込み
func writeHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm() // フォームを解析
	var log Log
	log.Name = request.Form["name"][0]
	log.Body = request.Form["body"][0]
	if log.Name == "" {
		log.Name = "名無し"
	}
	logs := loadLogs() // 既存データを読み出し
	log.ID = len(logs) + 1
	log.CTime = time.Now().Unix()
	logs = append(logs, log)                              // 追記
	saveLogs(logs)                                        // 保存
	http.Redirect(writer, request, "/", http.StatusFound) // リダイレクト
}

// 書き込みフォームを返す
func getForm() string {
	return "<div><form action='/write' method='POST'>" +
		"名前: <input type='text' name='name'><br>" +
		"本文: <input type='text' name='body' style='width:30em; '><br>" +
		"<input type='submit' value='書込'>" +
		"</form></div><hr>"
}

func main() {
	fmt.Println("server - http://localhost:8888")
	// URIに対応するハンドラを登録
	http.HandleFunc("/", showHandler)
	http.HandleFunc("/write", writeHandler)
	// サーバを起動
	http.ListenAndServe(":8888", nil)
}
