# isucon11-final

https://github.com/isucon/isucon11-final

## 設定例

```yaml
problem:
  name: isucon11-final
  options:
    benchmarker-path: /home/isucon/bench/benchmarker
    benchmarker-dir: /home/isucon/bench # ベンチマーカーを実行するときのディレクトリ
```

## ベンチマーカーの用意

`benchmarker` ディレクトリで `make` を実行するとビルドできます。

Makefile の `assets` ターゲットが、フロントエンドのビルド成果物
（`PUBLIC_FILES_DIR`）から静的ファイルのチェックサムを生成し、Go のソースに
埋め込みます。競技環境に配置するフロントエンドと同じ成果物でビルドしないと、
静的ファイルの検証に失敗します。チェックサムはビルド時に埋め込まれるため、
実行時に追加のディレクトリは必要ありません。
