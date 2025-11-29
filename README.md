ozen 🍱

ozen は、AIチャット（ChatGPT, Claude, GitHub Copilotなど）にコードを渡す際の「お膳立て」をするCLIツールです。

指示書（プロンプト）と複数のソースコードを読み込み、AIが理解しやすい形式に整形してクリップボードにコピーします。WSL環境でのクリップボード転送にも対応しています。

特徴

自動整形: ファイル名とコード内容をMarkdownのコードブロック形式で結合します。

プロンプト自動検出: カレントディレクトリの prompt.md や .github/copilot-instructions.md などの指示書を自動で検出し、コードの先頭に付与します。

クリップボード連携: -clip フラグを付けるだけで、整形済みテキストをクリップボードに送信します（WSL / Linux対応）。

使い方

基本的な使い方

カレントディレクトリに prompt.md がある場合、自動的に読み込まれます。

# prompt.md（もしあれば）と指定したソースコードを結合して標準出力に表示
ozen main.go utils.go


デフォルト以外の指示書を使いたい場合は、-prompt オプションで指定可能です。

クリップボードにコピー（推奨）

-clip フラグを使用すると、出力を直接クリップボードにコピーします。

# 全てのPythonソースを対象にする例
ozen -clip *.py

# 特定のディレクトリ以下のファイルを対象にする例
ozen -clip src/*.py


出力イメージ

以下のような形式のテキストが生成されます。これをAIチャットにペーストするだけで、文脈の共有が完了します。

実行コマンド例:

# 自動でprompt.mdが読み込まれ、Pythonファイルと結合される
ozen example/*.py


生成されるテキスト:

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


インストール

Releases ページから、お使いの環境に合わせたバイナリをダウンロードしてください。

Linux / WSL (x86_64)

# バイナリのダウンロード（バージョンは適宜変更してください）
curl -L -o ozen [https://github.com/Pikka2048/ozen/releases/download/v0.0.1/ozen-linux-amd64](https://github.com/Pikka2048/ozen/releases/download/v0.0.1/ozen-linux-amd64)

# 実行権限の付与
chmod +x ozen

# パスが通った場所へ移動
sudo mv ozen /usr/local/bin/


依存関係 (Linuxのみ)

WSL以外の純粋なLinux環境で使用する場合、クリップボード操作のために xclip が必要です。

sudo apt install xclip


※ WSL環境ではWindows側の clip.exe を経由するため、追加のインストールは不要です。
