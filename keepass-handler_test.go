package main

import (
	"reflect"
	"testing"

	"github.com/tobischo/gokeepasslib/v3"
)

var Handler *KeePassHandler = &KeePassHandler{}

// do not change the following 2 tests order!

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *KeePassHandler
	}{
		{
			name: "New Function",
			want: &KeePassHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := New()
			entry := gokeepasslib.NewEntry()
			entry.Values = append(entry.Values, mkValue("Title", "My GMail password"))
			entry.Values = append(entry.Values, mkValue("UserName", "example@gmail.com"))
			entry.Values = append(entry.Values, mkProtectedValue("Password", "hunter2"))

			k.db.Content.Root.Groups[0].Entries = append(k.db.Content.Root.Groups[0].Entries, entry)
			
			k.lockDB()

			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
