# go-redis-service-example

## 要件

- アクセスしたら jwt を発行し、redis にセッションを保持する
- 認証が通った場合、そうじゃない場合でレスポンスを変える

## エンドポイント

- signin
  - key を受け取って token を cookie にセット
  - redis に key-value を登録
- signiout
  - cookie から token 削除
  - redis から key-value を削除
- home
  - token を検証。redis をチェック。
  - 認証が通れば hello をレスポンス
  - 通らなければ 401 を返却

## docker-compose

- app サービス
  - 8080
- redis サービス
  - 6379
