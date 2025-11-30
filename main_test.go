package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "test.txt")

	// 作成前
	if fileExists(targetFile) {
		t.Errorf("Expected file %s to not exist", targetFile)
	}

	// 作成後
	if err := os.WriteFile(targetFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	if !fileExists(targetFile) {
		t.Errorf("Expected file %s to exist", targetFile)
	}
}

func TestGenerateTree(t *testing.T) {
	// テスト用ディレクトリ構造の作成
	// root/
	// ├── file1.txt
	// ├── .git/        (デフォルトで無視されるべきディレクトリ)
	// │   └── config
	// ├── sub/
	// │   └── file2.txt
	// └── node_modules/ (カスタム設定で無視するディレクトリ)
	//     └── lib.js

	root := t.TempDir()

	// ディレクトリ作成
	dirs := []string{".git", "sub", "node_modules"}
	for _, d := range dirs {
		if err := os.Mkdir(filepath.Join(root, d), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// ファイル作成
	files := []string{
		"file1.txt",
		filepath.Join(".git", "config"),
		filepath.Join("sub", "file2.txt"),
		filepath.Join("node_modules", "lib.js"),
	}
	for _, f := range files {
		path := filepath.Join(root, f)
		if err := os.WriteFile(path, []byte("dummy"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name          string
		config        TreeConfig
		shouldContain []string
		shouldExclude []string
	}{
		{
			name: "Normal tree with ignores",
			config: TreeConfig{
				MaxDepth: -1,
				Ignores:  map[string]bool{".git": true, "node_modules": true},
			},
			shouldContain: []string{
				"file1.txt",
				"sub",
				"file2.txt",
			},
			shouldExclude: []string{
				".git",
				"node_modules",
				"config", // .gitの中身
				"lib.js", // node_modulesの中身
			},
		},
		{
			name: "MaxDepth restriction (Depth=1)",
			config: TreeConfig{
				MaxDepth: 1, // ルート直下のみ表示、サブディレクトリの中身は表示しない
				Ignores:  map[string]bool{".git": true, "node_modules": true},
			},
			shouldContain: []string{
				"file1.txt",
				"sub",
			},
			shouldExclude: []string{
				"file2.txt", // subの中身なので表示されないはず
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTree(root, tt.config)
			if err != nil {
				t.Fatalf("GenerateTree() error = %v", err)
			}

			for _, s := range tt.shouldContain {
				if !strings.Contains(got, s) {
					t.Errorf("Output missing expected string: %s", s)
				}
			}

			for _, s := range tt.shouldExclude {
				if strings.Contains(got, s) {
					t.Errorf("Output contains excluded string: %s", s)
				}
			}
		})
	}
}
