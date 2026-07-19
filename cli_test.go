package litecli

import (
	"os"
	"reflect"
	"testing"
)

// TestExecute centralise tous les scénarios de parsing possibles pour valider le moteur.
func TestExecute(t *testing.T) {
	// On sauvegarde les vrais os.Args pour les restaurer à la fin du test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Définition de la structure pour chaque cas de test
	type testCase struct {
		name        string
		inputArgs   []string
		defString   string
		defInt      int
		defUint8    uint8
		wantString  string
		wantInt     int
		wantUint8   uint8
		wantPosArgs []string
	}

	tests := []testCase{
		{
			name:        "1. Valeurs par défaut quand aucun flag n'est passé",
			inputArgs:   []string{"gostamp", "image.jpg"},
			defString:   "./dist",
			defInt:      42,
			defUint8:    80,
			wantString:  "./dist",
			wantInt:     42,
			wantUint8:   80,
			wantPosArgs: []string{"image.jpg"},
		},
		{
			name:        "2. Flags longs collés avec '='",
			inputArgs:   []string{"gostamp", "--output=/tmp", "--quality=95", "pic.png"},
			defString:   "./dist",
			defInt:      0,
			defUint8:    80,
			wantString:  "/tmp",
			wantInt:     0,
			wantUint8:   95,
			wantPosArgs: []string{"pic.png"},
		},
		{
			name:        "3. Flags longs séparés par un espace",
			inputArgs:   []string{"gostamp", "--output", "/var/log", "--threads", "4", "file.tiff"},
			defString:   "./dist",
			defInt:      1,
			defUint8:    0,
			wantString:  "/var/log",
			wantInt:     4,
			wantUint8:   0,
			wantPosArgs: []string{"file.tiff"},
		},
		{
			name:        "4. Flags courts collés avec '='",
			inputArgs:   []string{"gostamp", "-o=/custom", "-t=8", "-q=100"},
			defString:   "",
			defInt:      0,
			defUint8:    0,
			wantString:  "/custom",
			wantInt:     8,
			wantUint8:   100,
			wantPosArgs: nil, // Aucun argument de position attendu
		},
		{
			name:        "5. Flags courts séparés par un espace",
			inputArgs:   []string{"gostamp", "-o", "bin", "-q", "50", "img1.jpg", "img2.jpg"},
			defString:   "",
			defInt:      0,
			defUint8:    0,
			wantString:  "bin",
			wantInt:     0,
			wantUint8:   50,
			wantPosArgs: []string{"img1.jpg", "img2.jpg"}, // Multiples arguments positionnels
		},
		{
			name:        "6. Flags inconnus (doivent être ignorés sans crash)",
			inputArgs:   []string{"gostamp", "--unknown-flag", "-x", "photo.jpg"},
			defString:   "default",
			defInt:      10,
			defUint8:    20,
			wantString:  "default",
			wantInt:     10,
			wantUint8:   20,
			wantPosArgs: []string{"photo.jpg"},
		},
		{
			name:        "7. Mélange complexe de formats et de positions",
			inputArgs:   []string{"gostamp", "first.jpg", "-o", "out", "--threads=12", "second.jpg", "-q", "75"},
			defString:   "",
			defInt:      1,
			defUint8:    0,
			wantString:  "out",
			wantInt:     12,
			wantUint8:   75,
			wantPosArgs: []string{"first.jpg", "second.jpg"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// On injecte les arguments simulés dans os.Args
			os.Args = tc.inputArgs

			// Variables cibles pour le test
			var gotString string
			var gotInt int
			var gotUint8 uint8
			var capturedPosArgs []string

			// Instanciation de la commande
			cmd := &Command{
				Short: "Test command",
				Run: func(c *Command, args []string) {
					capturedPosArgs = args
				},
			}

			// Enregistrement des flags (ce qui applique aussi les valeurs par défaut)
			cmd.StringVarP(&gotString, "output", 'o', tc.defString)
			cmd.IntVarP(&gotInt, "threads", 't', tc.defInt)
			cmd.Uint8VarP(&gotUint8, "quality", 'q', tc.defUint8)

			// Exécution du parsing
			cmd.Execute()

			// Vérifications des résultats
			if gotString != tc.wantString {
				t.Errorf("StringVarP échoué: attendu %q, obtenu %q", tc.wantString, gotString)
			}

			if gotInt != tc.wantInt {
				t.Errorf("IntVarP échoué: attendu %d, obtenu %d", tc.wantInt, gotInt)
			}

			if gotUint8 != tc.wantUint8 {
				t.Errorf("Uint8VarP échoué: attendu %d, obtenu %d", tc.wantUint8, gotUint8)
			}

			// Comparaison des slices d'arguments de position nettoyés
			if len(capturedPosArgs) == 0 && len(tc.wantPosArgs) == 0 {
				return // Slices vides considérés comme égaux
			}
			if !reflect.DeepEqual(capturedPosArgs, tc.wantPosArgs) {
				t.Errorf("Arguments de position incorrects: attendu %v, obtenu %v", tc.wantPosArgs, capturedPosArgs)
			}
		})
	}
}
