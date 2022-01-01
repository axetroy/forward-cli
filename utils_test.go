package forward

import "testing"

func Test_replaceHost(t *testing.T) {
	type args struct {
		content       string
		oldHost       string
		newHost       string
		proxyExternal bool
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
				oldHost: "example.com",
				newHost: "localhost:8080",
			},
			want: "http://localhost:8080",
		},
		{
			name: "1.1",
			args: args{
				content: "http://192.168.0.7:8443/webfolder/index.action",
				oldHost: "192.168.0.7:8443",
				newHost: "192.168.4.22:80",
			},
			want: "http://192.168.4.22:80/webfolder/index.action",
		},
		{
			name: "2",
			args: args{
				content:       "https://example.com.hk",
				oldHost:       "example.com",
				newHost:       "localhost:8080",
				proxyExternal: true,
			},
			want: "http://localhost:8080/?forward_url=https%3A%2F%2Fexample.com.hk",
		},
		{
			name: "3",
			args: args{
				content: "https://example.com/demo",
				oldHost: "example.com",
				newHost: "localhost:8080",
			},
			want: "http://localhost:8080/demo",
		},
		{
			name: "4",
			args: args{
				content:       "https://example.com.hk/demo",
				oldHost:       "example.com",
				newHost:       "localhost:8080",
				proxyExternal: true,
			},
			want: "http://localhost:8080/?forward_url=https%3A%2F%2Fexample.com.hk%2Fdemo",
		},
		{
			name: "5",
			args: args{
				content: "//example.com/demo",
				oldHost: "example.com",
				newHost: "localhost:8080",
			},
			want: "//localhost:8080/demo",
		},
		{
			name: "6",
			args: args{
				content:       "//www.baidu.com/s?wd=&%E7%99%BE%E5%BA%A6%E7%83%AD%E6%90%9C&sa=&ire_dl_gh_logo_texing&rsv_dl=&igh_logo_pcs",
				oldHost:       "www.baidu.com",
				newHost:       "localhost:8080",
				proxyExternal: true,
			},
			want: "//localhost:8080/s?wd=&%E7%99%BE%E5%BA%A6%E7%83%AD%E6%90%9C&sa=&ire_dl_gh_logo_texing&rsv_dl=&igh_logo_pcs",
		},
		{
			name: "7",
			args: args{
				content:       "https://passport.baidu.com/v2/?login&tpl=mn&u=http%3A%2F%2Fwww.baidu.com%2F",
				oldHost:       "www.baidu.com",
				newHost:       "localhost:8080",
				proxyExternal: true,
			},
			want: "http://localhost:8080/?forward_url=https%3A%2F%2Fpassport.baidu.com%2Fv2%2F%3Flogin%26tpl%3Dmn%26u%3Dhttp%253A%252F%252Flocalhost%253A8080%252F",
		},
		{
			name: "8",
			args: args{
				content: "file:///path/to/file",
				oldHost: "example.com",
				newHost: "localhost:8080",
			},
			want: "file:///path/to/file",
		},
		{
			name: "9",
			args: args{
				content: `E=/[^?#]*\//,S=/\/\.\//g,J=/\/[^/]+\/\.\.\//,U=/^([^/:]+)(\/.+)$/,V=/{([^{]+)}/g,R=/^\/\/.|:\//,T=/^.*?\/\/.*?\//`,
				oldHost: "example.com",
				newHost: "localhost:8080",
			},
			want: `E=/[^?#]*\//,S=/\/\.\//g,J=/\/[^/]+\/\.\.\//,U=/^([^/:]+)(\/.+)$/,V=/{([^{]+)}/g,R=/^\/\/.|:\//,T=/^.*?\/\/.*?\//`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceHost(tt.args.content, tt.args.oldHost, tt.args.newHost, false, tt.args.proxyExternal, []string{}); got != tt.want {
				t.Errorf("replaceHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
