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
`task run-server` で portal server が起動します。
