# proto

ポータルとベンチマーカーを実行する「runner」の通信を担当する。gRPCを使用する。

## 環境

protobufのコード生成を行うためには、[`buf`](https://buf.build/docs) が必要。[公式のインストール方法のドキュメント](https://buf.build/docs/installation/)を参照すること。
コード生成を行わない場合、bufは不要。

## コード生成

```sh
task proto:gen
```

もしくは

```sh
buf generate
```

## 処理の流れ

1. runnerがポータルにベンチマーカーキューの先頭を問い合わせる。ポータルは先頭のベンチマークIDを返す。(`GetBenchmarkJob`)
   1. キューが空の場合は、一定時間待機して再度問い合わせる。
2. runnerはベンチマークを実行し、client streamを使って一定時間ごとに結果をポータルに送信する。(`SendBenchmarkResult`)
3. ベンチマークが終了したら、runnerはポータルにベンチマークが終了したことを通知する。(`PostJobFinished`)
4. 1に戻る。
