# ansible

AWS 上に PISCON Portal と Runner を構築するための Ansible Playbook です。
`ubuntu` ユーザーで SSH 接続できる AWS の EC2 インスタンスをあらかじめ用意しておく必要があります。

## 環境準備

`uv` を使います。

```bash
uv sync
```

## ディレクトリ構成

```txt
.
├── README.md
├── ansible.cfg
├── group_vars
│   └── all
│       ├── local.yml # 実行ごとに値が異なるため
│       ├── secrets.yml # トークンなどの機密情報を含むファイル
│       └── vars.yml # 全ての変数定義。 local.yml と secrets.yml を参照する
├── hosts # インベントリファイル
├── portal.yml # ポータルサーバー構築用 Playbook
├── runner.yml # Runenr 構築用 Playbook
├── pyproject.toml
├── roles
│   ├── caddy # Caddy デプロイ用 role
│   ├── common # 共通設定用 role
│   ├── docker # Docker インストール用 role
│   ├── go # Go インストール用 role
│   ├── portal # ポータル起動用 role
│   └── runner # Runner 起動用 role
└── uv.lock
```

各 role の詳細については、各ディレクトリの `README.md` を参照してください。

## ファイルの準備

以下の3つのファイルを配置します。

### `group_vars/all/secrets.yml`

```yaml:secrets.yml
secrets:
  traq_client:
    id: "{{ traQのclient id }}"
    secret: "{{ traQのclient secret }}"
```

### `group_vars/all/local.yml`

```yaml:local.yml
local:
  problem:
    name: "{{ 問題名 (vars.yml の runner.problem_volumes に含まれるキーのうちどれか 1 つ) }}"
    instance_limit: {{ インスタンス数の上限 }}
    image_id: "{{ 問題のAMI ID }}"
    instance_type: "{{ 問題のインスタンスタイプ }}"

  aws:
    region: "{{ AWSのリージョン }}"
    subnet_id: "{{ AWSのサブネットID }}"
    security_group_id: "{{ AWSのセキュリティグループID }}"
    key_pair_name: "{{ AWSのキーペア名 }}"

  runner:
    portal_address: "{{ ポータルの gRPC サーバーのアドレス }}"

  admin:
    user_id: "{{ 初期管理者の traQ user ID (UUID) }}"
    user_name: "{{ 初期管理者の traQ user name }}"
```

### `hosts`

```ini:hosts
[portal]
piscon.trap.jp # ポータルのドメイン

[runner]
{{ RunnerのパブリックIPアドレス }}
{{ 複数指定できる }}
```

## 実行

### portal のデプロイ

```bash
uv run ansible-playbook portal.yml
```

### runner のデプロイ

```bash
uv run ansible-playbook runner.yml
```
