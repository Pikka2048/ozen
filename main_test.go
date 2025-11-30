package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existFile := filepath.Join(tmpDir, "exist.txt")
	if err := os.WriteFile(existFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path string
		want bool
	}{
		{existFile, true},
		{filepath.Join(tmpDir, "none.txt"), false},
	}

	for _, tt := range tests {
		if got := fileExists(tt.path); got != tt.want {
			t.Errorf("fileExists(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestCLI(t *testing.T) {
	// バイナリのビルド（テスト実行ごとのクリーンな環境のため）
	binPath := filepath.Join(t.TempDir(), "ozen_test")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", binPath, "main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// テスト用ディレクトリとファイル作成
	workDir := t.TempDir()
	promptContent := "Test Prompt Content\n"
	codeContent := "print('hello')\n"

	createFile := func(name, content string) {
		path := filepath.Join(workDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	createFile("prompt.md", promptContent)
	createFile("src/main.py", codeContent)
	createFile("other.md", "Custom Prompt")

	tests := []struct {
		name       string
		args       []string
		wantExit   int
		wantOutput []string // 含まれているべき文字列
	}{
		{
			name:       "No args",
			args:       []string{},
			wantExit:   1,
			wantOutput: []string{"ERROR: Specify the input file"},
		},
		{
			name:       "Default prompt and one file",
			args:       []string{"src/main.py"},
			wantExit:   0,
			wantOutput: []string{promptContent, "src/main.py", codeContent},
		},
		{
			name:       "Custom prompt flag",
			args:       []string{"-prompt", "other.md", "src/main.py"},
			wantExit:   0,
			wantOutput: []string{"Custom Prompt", "src/main.py"},
		},
		{
			name:       "Wildcard expansion",
			args:       []string{"src/*.py"},
			wantExit:   0,
			wantOutput: []string{promptContent, "src/main.py"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			cmd.Dir = workDir
			out, err := cmd.CombinedOutput()

			// Exit code check
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					t.Fatalf("Command execution failed: %v", err)
				}
			}

			if exitCode != tt.wantExit {
				t.Errorf("Exit code = %d, want %d", exitCode, tt.wantExit)
			}

			// Output content check
			gotStr := string(out)
			for _, want := range tt.wantOutput {
				if !strings.Contains(gotStr, want) {
					t.Errorf("Output missing %q. Got:\n%s", want, gotStr)
				}
			}
		})
	}
}
