# ozen

**ozen** は、ChatGPT、Claude、GitHub CopilotなどのAIチャット（LLM）へのコード共有を効率化するためのCLIツールです。

指示書（プロンプト）と複数のソースコードを読み込み、AIが文脈を理解しやすい形式に整形して出力します。また、標準出力だけでなくクリップボードへの直接転送（WSL/Linux対応）もサポートしています。

## Usage
カレントディレクトリに対象のソースコードを指定して実行します。整形されたテキストが標準出力に表示されます。

```bash
ozen main.py utils.py
````

ワイルドカードも使用可能です。

```bash
ozen src/*.py
```

デフォルトの探索対象以外の指示書を使用したい場合は、`-prompt` オプションで指定可能です。

```bash
ozen -prompt custom_instructions.md main.py
```

`-clip` フラグを使用すると、出力をコンソールに表示せず、直接クリップボードにコピーします。ブラウザのチャット欄に即座にペーストできるため効率的です。

```
ozen -clip src/*.py`は
ozen *.py | xclip -selection clipboardと同等。
```


カレントディレクトリの `prompt.md` や `.github/copilot-instructions.md` などの指示書を自動的に検出し、コードの冒頭に付与します。

### Example

```bash
# 自動的に prompt.md が読み込まれる想定
ozen example/*.py
```

````markdown
以下のルールを遵守して作業に取り組むこと。
- 過剰なコメントをしない
- 余計なimportをしない

```example/sample1.py
import os

def main():
    print("hello world")

if __name__ == "__main__":
    main()
```

```example/sample2.py
import sys

def func1(arg1):
    print(arg1)

def sum(arg1, arg2):
    return arg1, arg2
```

````

## Install

[Releases](https://github.com/Pikka2048/ozen/releases) ページより、ご使用の環境に合わせたバイナリをダウンロードしてください。

### Linux / WSL (x86\_64)

```bash
# バイナリのダウンロード
curl -L -o ozen https://github.com/Pikka2048/ozen/releases/download/v0.0.1/ozen-linux-amd64

# 実行権限の付与
chmod +x ozen

# PATHの通ったディレクトリへ移動
sudo mv ozen /usr/local/bin/
```

### Dependency  (Linux Only)

WSL以外のネイティブなLinux環境で使用する場合、クリップボード操作のために `xclip` のインストールが必要です。

```bash
sudo apt install xclip
```

> Note
> WSL環境ではWindows側の `clip.exe` を経由して動作するため、`xclip` のインストールは不要です。
