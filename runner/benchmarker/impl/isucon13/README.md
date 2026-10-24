# isucon13

https://github.com/isucon/isucon13

## 設定例

```yaml
problem:
  name: isucon13
  options:
    benchmarker-path: /home/isucon/bench
    benchmarker-dir: /home/isucon
    target-url: https://pipe.u.isucon.local # (任意) デフォルトはmatsuu/aws-isuconのAMIのURL
    dns-port: 53 # (任意) 対象サーバーのDNSポート
    enable-ssl: true # (任意) デフォルトはtrue
```

## 接続先

ジョブで受け取った対象IPを `--nameserver` に渡します。このポータルは問題サーバーのprivate IPをジョブに指定します。
同じチームの稼働中の問題サーバーのprivate IPもジョブで受け取り、対象IP以外を `--webapp` に渡します。
これにより、対象サーバーのDNSが同じチームの別サーバーのIPを返す複数台構成でも実行できます。
`target-url` はIPではなく、対象DNSが名前解決できるドメインを指定してください。
DNSが返すIPはジョブに含まれるチームの問題サーバーのprivate IPのいずれかと一致する必要があり、RunnerからそのIPのwebappへ到達できる必要があります。
Runnerから対象へのUDP/53とTCP/443の通信を許可してください。HTTPの場合は指定したHTTPポートも必要です。

[matsuu/aws-isucon の ISUCON13 AMI](https://github.com/matsuu/aws-isucon/tree/main/isucon13) は
ドメインを `u.isucon.local` に変更し、自己署名証明書を許可したベンチマーカーを配置しています。
Runnerのデフォルト値とAnsibleのテンプレートはこのAMIを前提に、`https://pipe.u.isucon.local` を使います。
AMIの初期状態でDNSが返す `127.0.0.1` は、問題サーバー側のDNSレコードで接続先webappのprivate IPへ変更してください。

## 結果の取得

公式ベンチマーカーの `run` コマンドを実行し、競技者向けの標準出力はINFOログも含めてそのままポータルに送ります。
運営向けの標準エラーはINFOログだけを除外し、WARN・ERRORなどの出力を残します。
ベンチマーカーが一時ファイルに書くログにはINFOも残ります。
スコアと合否は `--result-path` のJSONファイルから読み取ります。
途中経過のスコアは結果ファイルができるまで0です。
初期化・整合性チェックの失敗と、終了コード1を伴う最終チェックの失敗は競技上の失敗として扱います。
結果ファイルの欠落・不正や、それ以外の異常終了はRunnerのエラーとして扱います。

結果とログの保存先にはジョブごとの一時ディレクトリを使い、終了後に削除します。
