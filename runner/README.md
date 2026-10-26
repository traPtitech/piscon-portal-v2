# runner のログ

ログ全文は `logs/<benchmark ID>/stdout.gz` と `stderr.gz` に保存します。
Ansible の systemd サービスでは `/srv/runner/logs/`、Docker Compose では
`runner-logs` ボリュームの `/logs/` が保存先です。
`PISCON_LOG_DIR` 環境変数で保存先を変更できます。

`gzip -dc stdout.gz` / `zcat stderr.gz` で展開できます。
ログは自動削除しません。保存期間に応じて削除し、インスタンスやボリュームを
削除する前に、必要なログを退避してください。

runnerが複数台ある場合は、ベンチマークIDで保存先を探します。
Ansible構成では次のコマンドで、該当ログがあるホストを確認できます
（`BENCHMARK_ID`を対象IDに置き換えてください）。

```sh
ansible runner --become -m ansible.builtin.stat -a 'path=/srv/runner/logs/BENCHMARK_ID/stdout.gz'
```
