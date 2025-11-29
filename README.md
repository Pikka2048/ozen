# ozen 🍱

**ozen** は、AIチャット（ChatGPT, Claude, GitHub Copilotなど）にコードを渡すための「お膳立て」をするCLIツールです。

指示書（プロンプト）と複数のソースコードを読み込み、AIが理解しやすい形式に整形してクリップボードにコピーします。WSL環境でのクリップボード転送にも対応しています。

## 概要

  - **自動整形**: ファイル名とコード内容をMarkdownのコードブロック形式で結合します。
  - **プロンプト自動検出**: `prompt.md` や `instructions.md` などの指示書を自動で読み込み、コードの前に付与します。
  - **クリップボード連携**: `-clip` フラグをつけるだけで、整形済みテキストをクリップボードに送ります（WSL / Linux対応）。

## 使用例

### 基本的な使い方

カレントディレクトリの `prompt.md`（もしあれば）と、指定したソースコードを結合して標準出力に表示します。

```bash
ozen main.go utils.go

デフォルト（`prompt.md`, `.github/copilot-instructions.md` 等）以外の指示書を使いたい場合は `-prompt` で指定します。
```

### クリップボードにコピー（推奨）

`-clip` フラグを使用すると、出力を直接クリップボードにコピーします。

```bash
# 全てのPythonソースを対象
ozen -clip *.py

# 特定のディレクトリ以下のファイルを投げる
ozen -clip src/*.py
```

### 出力イメージ

以下のような形式のテキストが生成されます。目的の指示を書いたあとにペーストするだけで命令が完成します。

```bash
ozen example/*.py
```

````markdown
以下のルールを遵守して作業に取り組むこと。
- 過剰なコメントをしない
- 余計なimportをしない

[
```
example/sample1.py
import os

def main():
    print("hello world")

if __name__ == "__main__":
    main()
```

```
example/sample2.py
import sys

def func1(arg1):
    print(arg1)

def sum(arg1,arg2):
    return arg1, arg2
```
]

````

## インストール

[Releases](https://github.com/username/ozen/releases) ページから、お使いの環境に合わせたバイナリをダウンロードしてください。

### Linux / WSL (x86_64)

```bash
curl -L -o ozen https://github.com/Pikka2048/ozen/releases/download/v0.0.1/ozen-linux-amd64

# 実行権限の付与
chmod +x ozen

# パスが通った場所へ移動
sudo mv ozen /usr/local/bin/
````

### 依存関係 (Linuxのみ)

WSL以外の純粋なLinux環境で使用する場合、クリップボード操作のために `xclip` が必要です。

```bash
sudo apt install xclip
```

※ WSL環境では `clip.exe` を経由するため、追加のインストールは不要です。
