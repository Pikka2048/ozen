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

// WSLか確認
func isWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	cmd := exec.Command("uname", "-r")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	s := strings.ToLower(string(out))
	return strings.Contains(s, "microsoft") || strings.Contains(s, "wsl")
}

func main() {
	promptFile := flag.String("prompt", "", "any prompt file")
	useClip := flag.Bool("clip", false, "copy output to clipboard (auto-detects WSL/xclip)")
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
			str := fmt.Sprintf("\n```\n%s\n%s```\n", filename, string(data))
			contents = append(contents, str)
		}
	}

	result := fmt.Sprintf("%s\n%v", promptStr, contents)

	// クリップボードに渡す処理
	if *useClip {
		var cmd *exec.Cmd
		var msg string

		if isWSL() {
			// WSLならiconv -> clip.exe
			cmd = exec.Command("sh", "-c", "iconv -t cp932 | clip.exe")
			msg = "Copied to Windows clipboard (via WSL)."
		} else {
			// それ以外のLinuxならxclip
			cmd = exec.Command("xclip", "-selection", "clipboard")
			msg = "Copied to clipboard (via xclip)."
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
	} else {
		fmt.Println(result)
	}
}
