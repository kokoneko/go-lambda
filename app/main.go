package main

import (
	"encoding/json"
    "bytes"
    "log"
    "net/http"
    "os"
	"io/ioutil"

	event "app/event"
	resize "app/resize"
	"github.com/aws/aws-lambda-go/events"
)

var(
	runtimeApiEndpointPrefix string
)

func init() {
	// https://docs.aws.amazon.com/ja_jp/lambda/latest/dg/runtimes-api.html
	runtimeApiEndpointPrefix = "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API") + "/2018-06-01/runtime/invocation/"
}

func main() {
    log.Println("handler started")
	// イベントループ
	for {
        func() {
            // コンテキスト情報を取得
            resp, err := http.Get(runtimeApiEndpointPrefix + "next")
			if err != nil {
				log.Fatal(err)
			}
            defer func() {
                resp.Body.Close()
            }()

            // リクエストIDはヘッダに含まれる
            rId := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
            log.Printf("実行中のリクエストID" + rId)

			// 処理本体
			data, _ := ioutil.ReadAll(resp.Body)
			_, err = handle(data)
			if err != nil {
				http.Post(respErrorEndpoint(rId), "application/json", bytes.NewBuffer([]byte(rId)))
				log.Fatal(err)
			}

            // 最終的に真のランタイムに返すコンテンツはInvocation Response APIのリクエストボディに含める。
            http.Post(respEndpoint(rId), "application/json", bytes.NewBuffer([]byte(rId)))
        }()
    }
}

func handle(payload []byte) (string, error) {
	// https://github.com/aws/aws-lambda-go/blob/main/lambda/handler.go#L115 とかを参考に。
	var s3Event events.S3Event
	if err := json.Unmarshal(payload, &s3Event); err != nil {
		log.Fatal(err)
	}
	// ここからlambda処理の本体
	s3TrigerInfo := event.GetS3TrigerInfo(s3Event)
	resize.ExecResize(s3TrigerInfo.Bucket, s3TrigerInfo.Key)
	return "", nil
}

func respEndpoint(requestId string) string {
    return runtimeApiEndpointPrefix + requestId + "/response"
}

func respErrorEndpoint(requestId string) string {
    return runtimeApiEndpointPrefix + requestId + "/error"
}