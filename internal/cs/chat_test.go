package cs

import (
	"log"
	"testing"
	"time"
)

func TestNickGenerator(t *testing.T) {
	type args struct {
		nick string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test random generation",
			args: args{nick: "tomi"},
			want: "tomi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var prev string
			for i := 0; i < 100; i++ {
				nick := NickGenerator(tt.args.nick)
				if nick == prev {
					log.Println(nick)
					log.Println(prev)
					t.Fatal("the same number was generated")
				}
				log.Println(nick)
				prev = nick
				time.Sleep(1 * time.Millisecond)
			}

		})
	}
}
