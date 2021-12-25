package forward

import "testing"

func Test_replaceHost(t *testing.T) {
	type args struct {
		content string
		origin  string
		target  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				content: "https://example.com",
				origin:  "example.com",
				target:  "localhost:8080",
			},
			want: "http://localhost:8080",
		},
		{
			name: "2",
			args: args{
				content: "https://example.com.hk",
				origin:  "example.com",
				target:  "localhost:8080",
			},
			want: "https://example.com.hk",
		},
		{
			name: "3",
			args: args{
				content: "https://example.com/demo",
				origin:  "example.com",
				target:  "localhost:8080",
			},
			want: "http://localhost:8080/demo",
		},
		{
			name: "4",
			args: args{
				content: "https://example.com.hk/demo",
				origin:  "example.com",
				target:  "localhost:8080",
			},
			want: "https://example.com.hk/demo",
		},
		{
			name: "5",
			args: args{
				content: "//example.com/demo",
				origin:  "example.com",
				target:  "localhost:8080",
			},
			want: "//localhost:8080/demo",
		},
		{
			name: "6",
			args: args{
				content: "//www.baidu.com/s?wd=&%E7%99%BE%E5%BA%A6%E7%83%AD%E6%90%9C&sa=&ire_dl_gh_logo_texing&rsv_dl=&igh_logo_pcs",
				origin:  "www.baidu.com",
				target:  "localhost:8080",
			},
			want: "//localhost:8080/s?wd=&%E7%99%BE%E5%BA%A6%E7%83%AD%E6%90%9C&sa=&ire_dl_gh_logo_texing&rsv_dl=&igh_logo_pcs",
		},
		{
			name: "7",
			args: args{
				content: "https://passport.baidu.com/v2/?login&tpl=mn&u=http%3A%2F%2Fwww.baidu.com%2F",
				origin:  "www.baidu.com",
				target:  "localhost:8080",
			},
			want: "https://passport.baidu.com/v2/?login&tpl=mn&u=http%3A%2F%2Flocalhost%3A8080%2F",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceHost(tt.args.content, tt.args.origin, tt.args.target); got != tt.want {
				t.Errorf("replaceHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
