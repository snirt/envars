package cmd

// import (
// 	// "os"
// 	"testing"
// )

// func TestNew(t *testing.T) {
// 	// defer func() {
// 	// 	os.Remove(".env.kdbx")
// 	// }()
// 	tests := []struct {
// 		name string
// 		want *KeePassHandler
// 	}{
// 		{
// 			name: "New Function Write/Read",
// 			want: &KeePassHandler{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			k := New()
// 			k.AddVariables()
// 			k.lockDB()

// 			new_k := New()
// 			entries := new_k.db.Content.Root.Groups[0].Entries

// 			found := false
// 			for _, entry := range entries {
// 				for _, val := range entry.Values {
// 					if val.Key == "my_key" {
// 						found = true
// 					}
// 				}

// 			}
// 			if !found {
// 				t.Error("could not get the record from the decrypted db")
// 			}

// 		})
// 	}
// }
