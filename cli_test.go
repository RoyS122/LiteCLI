package litecli

import (
	"os"
	"reflect"
	"testing"
)

func TestAppRoutingAndParsing(t *testing.T) {

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	type testCase struct {
		name      string
		inputArgs []string

		rootDefaultQ   int
		wantRootCalled bool
		wantRootInt    int

		subDefaultOut   string
		subDefaultAlpha uint8
		wantSubCalled   bool
		wantSubStr      string
		wantSubUint8    uint8

		wantPosArgs []string
	}

	tests := []testCase{
		{
			name:           "1. Exécution Root par défaut (sans sous-commande)",
			inputArgs:      []string{"gostamp", "--quality=90", "root_file.jpg"},
			rootDefaultQ:   50,
			wantRootCalled: true,
			wantRootInt:    90,
			wantPosArgs:    []string{"root_file.jpg"},
		},
		{
			name:            "2. Routage vers sous-commande avec flags courts séparés",
			inputArgs:       []string{"gostamp", "process", "-o", "/tmp/dist", "-a", "120", "image1.png"},
			subDefaultOut:   "./dist",
			subDefaultAlpha: 75,
			wantSubCalled:   true,
			wantSubStr:      "/tmp/dist",
			wantSubUint8:    120,
			wantPosArgs:     []string{"image1.png"},
		},
		{
			name:            "3. Routage vers sous-commande avec flags longs collés (=)",
			inputArgs:       []string{"gostamp", "process", "--output=/var/out", "--watermark-alpha=200", "img2.webp", "img3.webp"},
			subDefaultOut:   "./dist",
			subDefaultAlpha: 75,
			wantSubCalled:   true,
			wantSubStr:      "/var/out",
			wantSubUint8:    200,
			wantPosArgs:     []string{"img2.webp", "img3.webp"},
		},
		{
			name:            "4. Sous-commande utilisant ses valeurs par défaut",
			inputArgs:       []string{"gostamp", "process", "default_test.jpg"},
			subDefaultOut:   "./fallback",
			subDefaultAlpha: 100,
			wantSubCalled:   true,
			wantSubStr:      "./fallback",
			wantSubUint8:    100,
			wantPosArgs:     []string{"default_test.jpg"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.inputArgs

			var rootCalled, subCalled bool
			var gotRootInt int
			var gotSubStr string
			var gotSubUint8 uint8
			var gotCapturedArgs []string

			rootCmd := &Command{
				Short: "Root command text",
				Run: func(cmd *Command, args []string) {
					rootCalled = true
					gotCapturedArgs = args
				},
			}
			rootCmd.IntVarP(&gotRootInt, "quality", 'q', tc.rootDefaultQ)

			processCmd := &Command{
				Use:   "process",
				Short: "Process subcommand text",
				Run: func(cmd *Command, args []string) {
					subCalled = true
					gotCapturedArgs = args
				},
			}
			processCmd.StringVarP(&gotSubStr, "output", 'o', tc.subDefaultOut)
			processCmd.Uint8VarP(&gotSubUint8, "watermark-alpha", 'a', tc.subDefaultAlpha)

			app := NewApp("gostamp", "Test App", "1.0.0", rootCmd)
			app.AddCommand(processCmd)

			app.Run()

			if rootCalled != tc.wantRootCalled {
				t.Errorf("Routage Root incorrect: attendu %v, obtenu %v", tc.wantRootCalled, rootCalled)
			}
			if subCalled != tc.wantSubCalled {
				t.Errorf("Routage Sous-Commande incorrect: attendu %v, obtenu %v", tc.wantSubCalled, subCalled)
			}

			if tc.wantRootCalled && gotRootInt != tc.wantRootInt {
				t.Errorf("Root IntVarP échoué: attendu %d, obtenu %d", tc.wantRootInt, gotRootInt)
			}

			if tc.wantSubCalled {
				if gotSubStr != tc.wantSubStr {
					t.Errorf("Sub StringVarP échoué: attendu %q, obtenu %q", tc.wantSubStr, gotSubStr)
				}
				if gotSubUint8 != tc.wantSubUint8 {
					t.Errorf("Sub Uint8VarP échoué: attendu %d, obtenu %d", tc.wantSubUint8, gotSubUint8)
				}
			}

			if len(gotCapturedArgs) == 0 && len(tc.wantPosArgs) == 0 {
				return
			}
			if !reflect.DeepEqual(gotCapturedArgs, tc.wantPosArgs) {
				t.Errorf("Arguments positionnels altérés: attendu %v, obtenu %v", tc.wantPosArgs, gotCapturedArgs)
			}
		})
	}
}
