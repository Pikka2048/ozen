package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// WSLか確認 (OS判定は呼び出し元で行う)
func isWSL() bool {
	cmd := exec.Command("uname", "-r")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	s := strings.ToLower(string(out))
	return strings.Contains(s, "microsoft") || strings.Contains(s, "wsl")
}

type TreeConfig struct {
	MaxDepth int // -1 で無制限
	Ignores  map[string]bool
}

func GenerateTree(rootPath string, config TreeConfig) (string, error) {
	var sb strings.Builder

	sb.WriteString(rootPath + "\n")

	// 再帰処理の開始
	err := appendTreeNodes(&sb, rootPath, "", 0, config)
	if err != nil {
		return sb.String(), err
	}

	return sb.String(), nil
}

// strings.Builder にツリー情報を再帰で書き込む内部関数
func appendTreeNodes(sb *strings.Builder, path string, prefix string, currentDepth int, config TreeConfig) error {
	// 深さ制限チェック
	if config.MaxDepth != -1 && currentDepth >= config.MaxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", path, err)
	}

	// フィルタリング処理
	var filtered []os.DirEntry
	for _, entry := range entries {
		if config.Ignores[entry.Name()] {
			continue
		}
		filtered = append(filtered, entry)
	}

	count := len(filtered)
	for i, entry := range filtered {
		isLast := i == count-1

		// 接続記号の決定
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		sb.WriteString(prefix)
		sb.WriteString(connector)
		sb.WriteString(entry.Name())
		sb.WriteString("\n")

		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			// 再帰呼び出し (sbのアドレスを渡す)
			err := appendTreeNodes(sb, filepath.Join(path, entry.Name()), newPrefix, currentDepth+1, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var defaultIgnores = []string{
	".git",
	".DS_Store",
}

func main() {
	promptFile := flag.String("prompt", "", "any prompt file")
	useClip := flag.Bool("clip", true, "copy output to clipboard (auto-detects WSL/xclip)")
	usePrint := flag.Bool("print", false, "copy output to clipboard (auto-detects WSL/xclip)")
	useTree := flag.Bool("tree", true, "tree command like directory listing (written in golang)")
	depthPtr := flag.Int("L", -1, "depth")
	ignorePtr := flag.String("ignore", "", "ignore file or directory name")
	flag.Parse()

	inputPatterns := flag.Args()

	if len(inputPatterns) == 0 {
		fmt.Println("ERROR: Specify the input file or pattern.")
		os.Exit(1)
	}

	// パターンに沿ったプロンプファイルが有ればそれを使う
	targetPrompt := *promptFile
	if targetPrompt == "" {
		defaults := []string{"prompt.md", ".github/copilot-instructions.md", "instructions.md"}
		for _, p := range defaults {
			if fileExists(p) {
				targetPrompt = p
				break
			}
		}
	}

	// プロンプトファイルを開く
	var promptStr string
	if targetPrompt != "" {
		prompt, err := os.ReadFile(targetPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s file open failed.\n", targetPrompt)
		} else {
			promptStr = string(prompt)
		}
	}

	// 読み込み処理
	var contents []string
	for _, pattern := range inputPatterns {
		filenames, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Pattern error %s\n", pattern)
			continue
		}

		for _, filename := range filenames {
			data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s file open failed.\n", filename)
				continue
			}
			str := fmt.Sprintf("\n```%s\n%s```\n", filename, string(data))
			contents = append(contents, str)
		}
	}

	// ディレクトリ構造取得
	var treeStr string
	if *useTree {
		ignores := make(map[string]bool)
		for _, name := range defaultIgnores {
			ignores[name] = true
		}
		// 除外ファイル
		if *ignorePtr != "" {
			parts := strings.Split(*ignorePtr, ",")
			for _, p := range parts {
				trimmed := strings.TrimSpace(p)
				if trimmed != "" {
					ignores[trimmed] = true
				}
			}
		}
		config := TreeConfig{
			MaxDepth: *depthPtr,
			Ignores:  ignores,
		}

		_treeStr, err := GenerateTree("./", config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		treeStr = _treeStr
	}

	// これまでの内容を合体
	content_str := strings.Join(contents, "")
	result := fmt.Sprintf("%s%v\n```\n%s\n```", promptStr, content_str, treeStr)

	// 以降は後処理
	// クリップボードに渡す処理
	if *useClip {
		var cmd *exec.Cmd
		var msg string

		if runtime.GOOS == "windows" {
			// Windowsネイティブ
			cmd = exec.Command("clip")
			msg = "Copied to Windows clipboard. Note: If you want to print, set -print option."
		} else if runtime.GOOS == "linux" && isWSL() {
			// WSLならiconv -> clip.exe
			cmd = exec.Command("sh", "-c", "iconv -t cp932 | clip.exe")
			msg = "Copied to Windows clipboard (via WSL). Note: If you want to print, set -print option."
		} else {
			// それ以外のLinuxならxclip
			cmd = exec.Command("xclip", "-selection", "clipboard")
			msg = "Copied to clipboard (via xclip). Note: If you want to print, set -print option."
		}

		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Pipe creation failed: %v\n", err)
			os.Exit(1)
		}

		go func() {
			defer stdin.Close()
			stdin.Write([]byte(result))
		}()

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Clipboard execution failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)
	}

	if *usePrint {
		fmt.Println(result)
	}
}
