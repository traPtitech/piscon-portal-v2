# piscon-portal-v2

## 開発環境
プロジェクトルートディレクトリに以下のように `.env` ファイルを用意してください。`CLIENT_ID`と`CLIENT_SECRET`はtraQのものです。
```
DB_USER=root
DB_PASSWORD=password
DB_HOST=db
DB_PORT=3306
DB_NAME=portal

ROOT_URL=http://localhost:8080

SESSION_SECRET=secret

CLIENT_ID=my-client-id
CLIENT_SECRET=my-client-secret
```
traQで作成するClientは リダイレクトURL に `http://localhost:8080/api/oauth2/callback`、スコープに `openid`, `profile` を指定します。
`task up` で portal server が 8080番 ポートで起動します。`task dev` でホットリロード付きで起動します。

Go 1.24で入った experimental な機能である `testing/synctest` を [runner/runner_test.go](runner/runner_test.go) で使っています。2025/2/26時点でこれを動かすためには、環境変数で `GOEXPERIMENT=synctest` を設定する必要があります。VSCode であれば、 `.vscode/settings.json` に

```json
{
  "go.toolsEnvVars": {
    "GOEXPERIMENT": "synctest"
  }
}
```

を書き込むことで環境変数を設定できます。

## protobuf について

[proto/README.md](proto/README.md) を参照してください。
