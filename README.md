# piscon-portal-v2

## 開発環境
プロジェクトルートディレクトリに以下のように `.env` ファイルを用意してください。`CLIENT_ID`と`CLIENT_SECRET`はtraQのものです。
```
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=portal

ROOT_URL=http://localhost:8080

SESSION_SECRET=secret

CLIENT_ID=my-client-id
CLIENT_SECRET=my-client-secret
```
traQで作成するClientは リダイレクトURL に `http://localhost:8080/api/oauth2/callback`、スコープに `openid`, `profile` を指定します。
`task run-server` で portal server が 8080番 ポートで起動します。

## protobuf について

[proto/README.md](proto/README.md) を参照してください。
