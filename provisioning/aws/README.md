# AWS CDK

AWS 上に PISCON Portal と Runner を構築するための AWS CDK です。

## 事前準備

デプロイには以下のツールが必要です。

- [Node.js](https://nodejs.org/ja)
- [AWS CLI](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/getting-started-install.html)
- [AWS CDK CLI](https://docs.aws.amazon.com/ja_jp/cdk/v2/guide/getting-started.html)

AWS への認証情報を設定しておく必要があります。詳しくは AWS CLI のドキュメントを参照してください。

## 環境準備

### 1. 依存パッケージのインストール

```bash
npm install
```

### 2. 環境変数の設定

プロジェクトルートに `.env` ファイルを作成し、以下の内容を記述します。

```dotenv:.env
# Runner インスタンスの AMI ID (必須)
RUNNER_AMI_ID="ami-xxxxxxxxxxxxxxxxx"

# 作成する Runner インスタンスの数 (任意、デフォルト: 1)
# RUNNER_COUNT=1

# Runner インスタンスのタイプ (任意、デフォルト: t3a.small)
# RUNNER_TYPE="t3a.small"

# EC2 に登録する SSH 公開鍵のパス (任意、デフォルト: ~/.ssh/id_ed25519.pub)
# SSH_PUBLIC_KEY="~/.ssh/id_rsa.pub"
```

## デプロイ

### 1. CDK のブートストラップ

初めて AWS 環境に CDK でデプロイする場合、以下のコマンドを実行してブートストラップ処理を行う必要があります。

```bash
cdk bootstrap
```

### 2. デプロイの実行

デプロイを実行します。

```bash
cdk deploy
```

デプロイが完了すると、Portal サーバーの IP アドレスなどがコンソールに出力されます。
