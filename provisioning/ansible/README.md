# ansible

## ファイルの準備

以下のファイルを配置します。

### `group_vars/all/secrets.yml`

```yaml:secrets.yml
secrets:
  traq_client:
    id: "{{ traQのclient id }}"
    secret: "{{ traQのclient secret }}"
```

### `group_vars/all/local.yml`

```yaml:local.yml
ignored:
  problem:
    instance_limit: {{ インスタンス数の上限 }}
    image_id: "{{ 問題のAMI ID }}"
    instance_type: "{{ 問題のインスタンスタイプ }}"

  aws:
    region: "{{ AWSのリージョン }}"
    subnet_id: "{{ AWSのサブネットID }}"
    security_group_id: "{{ AWSのセキュリティグループID }}"
    key_pair_name: "{{ AWSのキーペア名 }}"
```

### `hosts`

```ini:hosts
[portal]
piscon.trap.jp # ポータルのドメイン
```
