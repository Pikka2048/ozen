# ozen

**ozen** は、ChatGPT、Claude、GitHub CopilotなどのAIチャット（LLM）へのコード共有を効率化するためのCLIツールです。

プロンプトと複数のソースコードを読み込み、AIが文脈を理解しやすい形式に整形してクリップボードに出力します。(clip.exeかxclipを利用）

## Usage

```
Usage:
  ozen [patterns] [flags]

Flags:
      --clip             copy output to clipboard (auto-detects WSL/xclip) (default true)
  -L, --depth int        directory tree depth (default -1)
  -h, --help             help for ozen
      --ignore strings   ignore file or directory name
      --print            print output to stdout instead of clipboard
  -p, --prompt string    any prompt file
  -t, --tree             tree command like directory listing (default true)
```

 `prompt.md` や `.github/copilot-instructions.md` は自動的に検出し、`--prompt`としてセットされます。
 
### Example

```
$ ozen example/*.py
Copied to Windows clipboard (via WSL). Note: If you want to print, set --print option.
```

````
$ ozen example/*.py --print -t -L 2 --ignore .github (--printで標準出力に）

以下のルールを遵守して作業に取り組むこと。
- 過剰なコメントをしない
- 余計なimportをしない

以下に必要な背景情報(context)を提供します。

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

def sum(arg1,arg2):
    return arg1, arg2
```

Overview of the current directory.
```
./
├── .gitignore
├── Makefile
├── README.md
├── dist
│   ├── ozen-linux-amd64
│   └── ozen-windows-amd64.exe
├── example
│   ├── sample1.py
│   └── sample2.py
├── go.mod
├── go.sum
├── main.go
├── main_test.go
└── prompt.md
```
````


## Install

[Releases](https://github.com/Pikka2048/ozen/releases) ページより、ご使用の環境に合わせたバイナリをダウンロードしてください。

### Linux / WSL / Windows (x86\_64)

```bash
# バイナリのダウンロード
curl -L -o ozen https://github.com/Pikka2048/ozen/releases/download/v0.0.1/ozen-linux-amd64

# 実行権限の付与
chmod +x ozen

# PATHの通ったディレクトリへ移動
sudo mv ozen /usr/local/bin/
```
