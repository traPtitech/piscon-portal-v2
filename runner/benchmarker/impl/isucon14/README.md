# isucon14

https://github.com/isucon/isucon14

## 設定例

```yaml
problem:
  name: isucon14
  options:
    benchmarker-path: /home/isucon/bench
    benchmarker-ip: 192.168.0.3 # ベンチマーカー自身のipアドレス。ターゲットから決済サーバーとして到達できる必要がある
    benchmarker-dir: /home/isucon # ベンチマーカーを実行するときのディレクトリ
    payment-bind-port: 12346 # (任意) ベンチマーカーが立てる決済サーバーのポート。デフォルトは12346
    target-url: https://isuride.xiv.isucon.net # (任意) HostヘッダーとSNIに使うURL。デフォルトは https://isuride.xiv.isucon.net
    skip-static-sanity-check: false # (任意) 静的ファイルの検証をスキップする。デフォルトはfalse
```

## 注意

- ターゲットへは `<ターゲットのIP>:443` に接続し、Host ヘッダーと SNI には `target-url` のホストを使います。
  競技環境の nginx は `*.xiv.isucon.net` 宛てのリクエストにしかアプリを返しません。
  証明書の検証はベンチマーカー側で行われないため、自己署名証明書でも動きます。
- ベンチマーカーは決済サーバーのモックを `payment-bind-port` で立て、ターゲットはそこにリクエストを送ります。
  [matsuu/aws-isucon](https://github.com/matsuu/aws-isucon/tree/main/isucon14) の AMI では
  `isuride-payment_mock` が 12345 番を使っているため、デフォルトを 12346 にしています。
- ベンチマーカーは静的ファイルのハッシュを埋め込んでおり、競技環境のフロントエンドと一致しない場合は失敗します。
- Prepare (initialize や静的ファイルの検証など) で失敗した場合、ベンチマーカーは最終結果を送らずに終了するため、
  スコア 0 の失敗として扱います。
