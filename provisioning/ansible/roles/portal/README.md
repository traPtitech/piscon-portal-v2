# portal

PISCON Portal v2 をインストールし、サービスとして起動します。
docker compose を使ってアプリサーバー、フロントエンド、MariaDB を起動します。これらは systemd のサービスとして起動します。
アプリサーバーとフロントエンドは、GitHub Container Registry からコンテナイメージを取得します。
