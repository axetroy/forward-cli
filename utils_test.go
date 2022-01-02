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
		{
			name: "10",
			args: args{
				content: `[\x22gws-wiz\x22,\x22\x22,\x22\x22,\x22\x22,null,1,0,0,11,\x22zh-CN\x22,\x22OkGroE8_h_-s0Ft8ObIJjo8cYg8\x22,\x222e63be090c64c5eb35509377f34f3609587d309f\x22,\x22JEjRYaqIO5K0mAXX3pC4Bw\x22,0,\x22zh-CN\x22,null,null,null,3,5,null,8,null,\x22\x22,-1,0,0,null,1,0,null,0,0,0,1,0,0,8,-1,null,0,null,null,1,0,0,0,0,0.1,null,0,100,0,null,1.15,1,null,null,null,0,null,0,0,0,6,0,null,null,null,null,null,0,0,0,0,null,null,0,null,null,0,0,0,null,null,null,null,null,null,null,0,null,0,0,0,null,\x22\x22,0,1,0,-1,null,0]`,
				oldHost: "www.google.com",
				newHost: "localhost",
			},
			want: `[\x22gws-wiz\x22,\x22\x22,\x22\x22,\x22\x22,null,1,0,0,11,\x22zh-CN\x22,\x22OkGroE8_h_-s0Ft8ObIJjo8cYg8\x22,\x222e63be090c64c5eb35509377f34f3609587d309f\x22,\x22JEjRYaqIO5K0mAXX3pC4Bw\x22,0,\x22zh-CN\x22,null,null,null,3,5,null,8,null,\x22\x22,-1,0,0,null,1,0,null,0,0,0,1,0,0,8,-1,null,0,null,null,1,0,0,0,0,0.1,null,0,100,0,null,1.15,1,null,null,null,0,null,0,0,0,6,0,null,null,null,null,null,0,0,0,0,null,null,0,null,null,0,0,0,null,null,null,null,null,null,null,0,null,0,0,0,null,\x22\x22,0,1,0,-1,null,0]`,
		},
		{
			name: "11",
			args: args{
				content: `window.jsl.dh=function(d,e,c){try{var f=document.getElementById(d);if(f)f.innerHTML=e,c&&c();else{var a={id:d,script:String(!!c),milestone:String(google.jslm||0)};google.jsla&&(a.async=google.jsla);var g=document.createElement("div");g.innerHTML=e;var b=g.children[0];b&&(a.tag=b.tagName,a["class"]=String(b.className||null),a.name=String(b.getAttribute("jsname")));google.ml(Error("Missing ID."),!1,a)}}catch(h){google.ml(h,!0,{"jsl.dh":!0})}};(function(){var x=true;google.jslm=x?2:1;})();google.x(null, function(){(function(){(function(){google.csct={};google.csct.ps='AOvVaw17ag9mz-2UL3tGGKniglcH\x26ust\x3d1641191845010685';})();})();(function(){(function(){google.csct.rw=true;})();})();(function(){(function(){google.csct.rl=true;})();})();(function(){google.drty&&google.drty(undefined,true);})();});google.drty&&google.drty(undefined,true);</script></body></html>`,
				oldHost: "www.google.com",
				newHost: "localhost",
			},
			want: `window.jsl.dh=function(d,e,c){try{var f=document.getElementById(d);if(f)f.innerHTML=e,c&&c();else{var a={id:d,script:String(!!c),milestone:String(google.jslm||0)};google.jsla&&(a.async=google.jsla);var g=document.createElement("div");g.innerHTML=e;var b=g.children[0];b&&(a.tag=b.tagName,a["class"]=String(b.className||null),a.name=String(b.getAttribute("jsname")));google.ml(Error("Missing ID."),!1,a)}}catch(h){google.ml(h,!0,{"jsl.dh":!0})}};(function(){var x=true;google.jslm=x?2:1;})();google.x(null, function(){(function(){(function(){google.csct={};google.csct.ps='AOvVaw17ag9mz-2UL3tGGKniglcH\x26ust\x3d1641191845010685';})();})();(function(){(function(){google.csct.rw=true;})();})();(function(){(function(){google.csct.rl=true;})();})();(function(){google.drty&&google.drty(undefined,true);})();});google.drty&&google.drty(undefined,true);</script></body></html>`,
		},
		{
			name: "12",
			args: args{
				content: `(function(){var c=Date.now();if(google.timers&&google.timers.load.t){for(var a=document.getElementsByTagName("img"),d=0,b=void 0;b=a[d++];)google.c.setup(b,!1,void 0);google.c.frt=!1;google.c.e("load","imn",String(a.length));google.c.ubr(!0,c);google.c.glu&&google.c.glu();google.rll(window,!1,function(){google.tick("load","ol");google.c.u("pr")})}})();}).call(this);(function(){google.jl={attn:false,blt:'none',chnk:0,dw:false,dwu:true,emtn:0,end:0,ine:false,lls:'default',pdt:0,rep:0,snet:true,strt:0,ubm:false,uwp:true};})();(function(){var pmc='{\x22aa\x22:{},\x22abd\x22:{\x22abd\x22:false,\x22deb\x22:false,\x22det\x22:false},\x22async\x22:{},\x22cdos\x22:{\x22cdobsel\x22:false},\x22cr\x22:{\x22qir\x22:false,\x22rctj\x22:true,\x22ref\x22:false,\x22uff\x22:false},\x22csi\x22:{},\x22d\x22:{},\x22dpf\x22:{},\x22dvl\x22:{\x22cookie_secure\x22:true,\x22cookie_timeout\x22:21600,\x22jsc\x22:\x22[null,null,null,30000,null,null,null,2,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,[\\\x2286400000\\\x22,\\\x22604800000\\\x22,2],null,null,21600000,null,null,1,null,null,null,null,null,1]\x22,\x22msg_err\x22:\x22无法获取位置信息\x22,\x22msg_gps\x22:\x22使用 GPS\x22,\x22msg_unk\x22:\x22未知\x22,\x22msg_upd\x22:\x22更新位置信息\x22,\x22msg_use\x22:\x22使用确切位置\x22,\x22use_local_storage_fallback\x22:false},\x22gf\x22:{\x22pid\x22:196,\x22si\x22:true},\x22hsm\x22:{},\x22jsa\x22:{\x22csi\x22:true,\x22csir\x22:100},\x22mu\x22:{\x22murl\x22:\x22https://adservice.google.com/adsid/google/ui\x22},\x22pHXghd\x22:{},\x22sb_wiz\x22:{\x22rfs\x22:[],\x22scq\x22:\x22\x22,\x22stok\x22:\x22OkGroE8_h_-s0Ft8ObIJjo8cYg8\x22,\x22ueh\x22:\x222e63be09_0c64c5eb_35509377_f34f3609_587d309f\x22},\x22sf\x22:{}}';google.pmc=JSON.parse(pmc);})();(function(){var r=['sb_wiz','aa','abd','async','dvl','mu','pHXghd','sf'];google.plm(r);})();(function(){var m=['BrDNik','[\x22gws-wiz\x22,\x22\x22,\x22\x22,\x22\x22,null,1,0,0,11,\x22zh-CN\x22,\x22OkGroE8_h_-s0Ft8ObIJjo8cYg8\x22,\x222e63be090c64c5eb35509377f34f3609587d309f\x22,\x22JEjRYaqIO5K0mAXX3pC4Bw\x22,0,\x22zh-CN\x22,null,null,null,3,5,null,8,null,\x22\x22,-1,0,0,null,1,0,null,0,0,0,1,0,0,8,-1,null,0,null,null,1,0,0,0,0,0.1,null,0,100,0,null,1.15,1,null,null,null,0,null,0,0,0,6,0,null,null,null,null,null,0,0,0,0,null,null,0,null,null,0,0,0,null,null,null,null,null,null,null,0,null,0,0,0,null,\x22\x22,0,1,0,-1,null,0]'];
var a=m;window.W_jd=window.W_jd||{};for(var b=0;b<a.length;b+=2)window.W_jd[a[b]]=JSON.parse(a[b+1]);})();(function(){window.WIZ_global_data={"GWsdKe":"zh-Hans-TW","SNlM0e":"AD7QlO4fYWlmotJfBMsAjUNJygCg:1641105445019","LVIXXb":"1","zChJod":"%.@.]","Yllh3e":"%.@.1641105444967722,178657810,1996762967]","w2btAe":"%.@.\"101199000968061242063\",\"101199000968061242063\",\"0\",null,null,null,1]","QrtxK":"0","eptZe":"/wizrpcui/_/WizRpcUi/","S06Grb":"101199000968061242063"};window.IJ_values=[false,true,true,true,false,false,false,24,"none",true,"1px 1px 15px 0px #171717",false,"rgba(255,255,255,.54)","rgba(255,255,255,.26)","#000","rgba(0,0,0,.3)",true,"none","#424548",true,false,false,true,false,"#609beb","#8ab4f8","1px 1px 15px 0px #171717",true,false,36,24,28,6,true,false,false,false,false,false,"#3c4043",10,true,false,false,"#202124","#e8eaed",false,"#303134","0px 5px 26px 0px rgba(0,0,0,0.5), 0px 20px 28px 0px rgba(0,0,0,0.5)","#4285f4",false,true,false,"#8ab4f8",false,true,false,false,"#fff","#4487f6","#48a1ff","#a4c2ff","#219540","#41b85f","#4eb66e","#ff7d70","#ff897e","#000","#1aa863","#212327","#050505","#0a0a0a","#111","#9aa0a6","#8a8a8a","#bdc1c6","#bdc1c6","#bdc1c6","#000","#8a4a00","#824300","#b85100","18px","#3c4043","#e8eaed","#e8eaed","#3c4043",14,"#e8eaed",40,"#e8eaed",false,"#8a8a8a","#ddd","#ff7a8e","#bdc1c6","arial,sans-serif-medium,sans-serif","arial,sans-serif","#bdc1c6","#292e36","#bdc1c6","#868b90","#48a1ff",false,false,false,false,false,false,true,false,false,false,"1px 1px 15px 0px #171717",false,false,"#3c4043","rgba(255,255,255,.26)","#9aa0a6","#bdc1c6","rgba(204,204,204,.15)","rgba(204,204,204,.25)","rgba(102,102,102,.2)","rgba(102,102,102,.4)","rgba(255,255,255,.12)","#3c4043","#fff","rgba(0,0,0,.3)","#000","#bdc1c6","#000","Roboto,RobotoDraft,Helvetica,Arial,sans-serif","14px","500","500","pointer","0 1px 1px rgba(0,0,0,.16)",true,24,"#000","1px 1px 15px 0px #171717","#dadce0",200,true,true,false,false,true,true,false,true,14,"#202124","#303134",false,"1px solid #3c4043","none","arial,sans-serif-medium,sans-serif","Google Sans,arial,sans-serif-medium,sans-serif","#3c4043","1px solid #3c4043","1px solid #5f6368","rgba(255,255,255,0.1)","#3c4043","#202124","#8ab4f8","#3c4043","#bdc1c6","#9aa0a6",false,true,true,false,false,false,false,false,false,false,false,false,true,false,false,false,true,false,false,false,false,false,false,false,"8px","#3c4043",false,true,false,"%.@.\"101199000968061242063\",\"101199000968061242063\",\"0\",null,null,null,1]","0","%.@.null,1,1,null,[null,757,1440]]","LdSGxtnn+iuvNIwJ+JnlLg\u003d\u003d","%.@.\"#424548\"]","%.@.0]","%.@.0]","%.@.\"0px 5px 26px 0px rgba(0,0,0,0.5),0px 20px 28px 0px rgba(0,0,0,0.5)\",\"#303134\"]","%.@.0,null,null,36,28,6,0.3,null,14,null,null,null,null,null,\"#bdc1c6\",\"#9aa0a6\",null,\"#bdc1c6\",null,null,null,null,null,null,\"#1a73e8\",\"#fabb05\",\"#fff\",\"#1a73e8\",\"#d1d1d1\",\"#fff\",null,null,null,14,500,\"#51a6ff\",null,\"#8ab4f8\",\"#303134\"]",null,"%.@.[],0,null,1,1]","zh-Hans-TW","%.@.\"13px\",\"16px\",\"11px\",13,16,11,\"8px\",8,20]","zh_Hans_TW","%.@.\"10px\",10,\"16px\",16,\"18px\"]","%.@.\"14px\",14]","%.@.40,32,14]",null,"%.@.\"1px 1px 15px 0px #171717\"]","%.@.0,\"14px\",\"500\",\"500\",\"0 1px 1px rgba(0,0,0,.16)\",\"pointer\",\"#fff\",\"rgba(255,255,255,.26)\",\"#9aa0a6\",\"#bdc1c6\",\"rgba(204,204,204,.15)\",\"rgba(204,204,204,.25)\",\"rgba(102,102,102,.2)\",\"rgba(102,102,102,.4)\",\"#1aa863\",\"#4487f6\",\"#a4c2ff\",\"#ff7d70\",\"#8a4a00\",\"#111\",\"#050505\",\"#bdc1c6\",\"#4f861f\",\"rgba(255,255,255,.12)\",null,\"#000\",\"rgba(0,0,0,.3)\",\"#000\",\"#bdc1c6\",\"#000\",null,0]","%.@.\"20px\",\"500\",\"400\",\"13px\",\"15px\",\"15px\",\"Roboto,RobotoDraft,Helvetica,Arial,sans-serif\",\"24px\",\"400\",\"32px\",\"24px\"]",false,"","%.@.null,null,null,null,\"20px\",\"20px\",\"18px\",\"40px\",\"36px\",\"32px\",null,null,null,null,null,null,\"#202124\",null,null,null,\"#202124\",null,null,null,\"rgba(138,180,248,0.24)\",null,\"rgba(138,180,248,0.24)\",null,null,\"16px\",\"12px\",\"8px\",\"4px\",\"#202124\",\"rgba(138,180,248,0.24)\",\"#d2e3fc\",\"transparent\",\"#8ab4f8\",\"#5f6368\",\"999rem\",\"8px\",\"#d2e3fc\",\"transparent\",\"#dadce0\",\"#5f6368\",\"#d2e3fc\",\"transparent\",\"#8ab4f8\",\"#5f6368\",\"999rem\",\"Google Sans,arial,sans-serif-medium,sans-serif\",\"20px\",\"14px\",\"500\"]","%.@.\"#bdc1c6\",\"#bdc1c6\",\"#8ab4f8\",null,\"#9aa0a6\",\"#8ab4f8\",\"#c58af9\",null,null,\"#202124\",\"#8ab4f8\",\"#202124\",\"#394457\",\"#d2e3fc\",\"#303134\",\"#bdc1c6\",\"#fff\",\"#3c4043\",\"#202124\",\"#fff\",\"#202124\",\"#fff\",\"#81c995\",\"#f28b82\",\"#fdd663\",\"#3c4043\",\"#202124\",\"rgba(0,0,0,0.6)\",\"#bdc1c6\",\"#3c4043\"]","%.@.null,\"none\",null,\"0px 1px 3px hsla(0,0%,9%,0.24)\",null,\"0px 2px 6px hsla(0,0%,9%,0.32)\",null,\"0px 4px 12px hsla(0,0%,9%,0.9)\",null,null,\"1px solid  #5f6368\",\"none\",\"none\",\"none\"]","%.@.\"Google Sans,arial,sans-serif\",\"Google Sans,arial,sans-serif-medium,sans-serif\",\"arial,sans-serif\",\"arial,sans-serif-medium,sans-serif\",\"arial,sans-serif-light,sans-serif\"]","%.@.\"16px\",\"12px\",\"0px\",\"8px\",\"4px\",\"2px\",\"20px\",\"24px\"]","%.@.\"#8ab4f8\",\"#8ab4f8\"]","%.@.null,null,null,null,null,null,null,\"12px\",\"8px\",\"4px\",\"16px\",\"2px\",\"999rem\",\"0px\"]","%.@.\"700\",\"400\",\"underline\",\"none\",\"capitalize\",\"none\",\"uppercase\",\"none\",\"500\",\"lowercase\",\"italic\",\"-1px\",\"0.3px\"]","%.@.\"20px\",\"26px\",\"400\",\"Google Sans,arial,sans-serif\",null,\"arial,sans-serif\",\"14px\",\"400\",\"22px\",null,\"16px\",\"24px\",\"400\",\"Google Sans,arial,sans-serif\",null,\"Google Sans,arial,sans-serif\",\"60px\",\"48px\",\"-1px\",null,\"400\",\"Google Sans,arial,sans-serif\",\"36px\",\"400\",\"48px\",null,\"Google Sans,arial,sans-serif\",\"36px\",\"28px\",null,\"400\",null,\"arial,sans-serif\",\"24px\",\"18px\",null,\"400\",\"arial,sans-serif\",\"16px\",\"12px\",null,\"400\",\"arial,sans-serif\",\"22px\",\"16px\",null,\"400\",\"arial,sans-serif\",\"26px\",\"20px\",null,\"400\",\"arial,sans-serif\",\"20px\",\"16px\",null,\"400\",\"arial,sans-serif\",\"18px\",\"14px\",null,\"400\",\"Google Sans,arial,sans-serif\",\"32px\",\"24px\",null,\"500\"]","%.@.4]","%.@.\"14px\",14,\"16px\",16,\"0\",0,\"none\",632,\"1px solid #3c4043\",\"normal\",\"normal\",\"#9aa0a6\",\"12px\",\"1.34\",\"1px solid #3c4043\",\"none\",\"0\",\"none\",\"none\",\"none\",\"none\",\"6px\"]","%.@.\"0\"]","%.@.\"rgba(0,0,0,0.0)\",\"rgba(0,0,0,0.54)\",\"rgba(0,0,0,0.8)\",\"rgba(248, 249, 250, 0.85)\",\"#202124\",\"#dadce0\",\"rgba(218, 220, 224, 0.0)\",\"rgba(218, 220, 224, 0.7)\",\"#dadce0\",\"#f8f9fa\",\"#000\",\"#1a73e8\",\"#dadce0\",\"#fff\",\"#fff\",\"#e8eaed\"]","%.@.\"#dddee1\",\"#868b90\",\"#bdc1c6\",\"#bcc0c3\",\"#000\",\"rgba(0,0,0,.7)\",28,24,26,20,16,-2,0,-4,2,0,0,24,20,20,14,12]","%.@.\"20px\",20,\"14px\",14,\"#e8eaed\"]","troy450409405@gmail.com",true,"101199000968061242063","%.@.1]"];})();(function(){google.llirm='400px';google.ldi={};google.pim={};})();
window.jsl=window.jsl||{};window.jsl.dh=function(d,e,c){try{var f=document.getElementById(d);if(f)f.innerHTML=e,c&&c();else{var a={id:d,script:String(!!c),milestone:String(google.jslm||0)};google.jsla&&(a.async=google.jsla);var g=document.createElement("div");g.innerHTML=e;var b=g.children[0];b&&(a.tag=b.tagName,a["class"]=String(b.className||null),a.name=String(b.getAttribute("jsname")));google.ml(Error("Missing ID."),!1,a)}}catch(h){google.ml(h,!0,{"jsl.dh":!0})}};(function(){var x=true;
google.jslm=x?2:1;})();google.x(null, function(){(function(){(function(){google.csct={};google.csct.ps='AOvVaw17ag9mz-2UL3tGGKniglcH\x26ust\x3d1641191845010685';})();})();(function(){(function(){google.csct.rw=true;})();})();(function(){(function(){google.csct.rl=true;})();})();(function(){google.drty&&google.drty(undefined,true);})();});google.drty&&google.drty(undefined,true);`,
				oldHost: "www.google.com",
				newHost: "localhost",
			},
			want: `(function(){var c=Date.now();if(google.timers&&google.timers.load.t){for(var a=document.getElementsByTagName("img"),d=0,b=void 0;b=a[d++];)google.c.setup(b,!1,void 0);google.c.frt=!1;google.c.e("load","imn",String(a.length));google.c.ubr(!0,c);google.c.glu&&google.c.glu();google.rll(window,!1,function(){google.tick("load","ol");google.c.u("pr")})}})();}).call(this);(function(){google.jl={attn:false,blt:'none',chnk:0,dw:false,dwu:true,emtn:0,end:0,ine:false,lls:'default',pdt:0,rep:0,snet:true,strt:0,ubm:false,uwp:true};})();(function(){var pmc='{\x22aa\x22:{},\x22abd\x22:{\x22abd\x22:false,\x22deb\x22:false,\x22det\x22:false},\x22async\x22:{},\x22cdos\x22:{\x22cdobsel\x22:false},\x22cr\x22:{\x22qir\x22:false,\x22rctj\x22:true,\x22ref\x22:false,\x22uff\x22:false},\x22csi\x22:{},\x22d\x22:{},\x22dpf\x22:{},\x22dvl\x22:{\x22cookie_secure\x22:true,\x22cookie_timeout\x22:21600,\x22jsc\x22:\x22[null,null,null,30000,null,null,null,2,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,[\\\x2286400000\\\x22,\\\x22604800000\\\x22,2],null,null,21600000,null,null,1,null,null,null,null,null,1]\x22,\x22msg_err\x22:\x22无法获取位置信息\x22,\x22msg_gps\x22:\x22使用 GPS\x22,\x22msg_unk\x22:\x22未知\x22,\x22msg_upd\x22:\x22更新位置信息\x22,\x22msg_use\x22:\x22使用确切位置\x22,\x22use_local_storage_fallback\x22:false},\x22gf\x22:{\x22pid\x22:196,\x22si\x22:true},\x22hsm\x22:{},\x22jsa\x22:{\x22csi\x22:true,\x22csir\x22:100},\x22mu\x22:{\x22murl\x22:\x22https://adservice.google.com/adsid/google/ui\x22},\x22pHXghd\x22:{},\x22sb_wiz\x22:{\x22rfs\x22:[],\x22scq\x22:\x22\x22,\x22stok\x22:\x22OkGroE8_h_-s0Ft8ObIJjo8cYg8\x22,\x22ueh\x22:\x222e63be09_0c64c5eb_35509377_f34f3609_587d309f\x22},\x22sf\x22:{}}';google.pmc=JSON.parse(pmc);})();(function(){var r=['sb_wiz','aa','abd','async','dvl','mu','pHXghd','sf'];google.plm(r);})();(function(){var m=['BrDNik','[\x22gws-wiz\x22,\x22\x22,\x22\x22,\x22\x22,null,1,0,0,11,\x22zh-CN\x22,\x22OkGroE8_h_-s0Ft8ObIJjo8cYg8\x22,\x222e63be090c64c5eb35509377f34f3609587d309f\x22,\x22JEjRYaqIO5K0mAXX3pC4Bw\x22,0,\x22zh-CN\x22,null,null,null,3,5,null,8,null,\x22\x22,-1,0,0,null,1,0,null,0,0,0,1,0,0,8,-1,null,0,null,null,1,0,0,0,0,0.1,null,0,100,0,null,1.15,1,null,null,null,0,null,0,0,0,6,0,null,null,null,null,null,0,0,0,0,null,null,0,null,null,0,0,0,null,null,null,null,null,null,null,0,null,0,0,0,null,\x22\x22,0,1,0,-1,null,0]'];
var a=m;window.W_jd=window.W_jd||{};for(var b=0;b<a.length;b+=2)window.W_jd[a[b]]=JSON.parse(a[b+1]);})();(function(){window.WIZ_global_data={"GWsdKe":"zh-Hans-TW","SNlM0e":"AD7QlO4fYWlmotJfBMsAjUNJygCg:1641105445019","LVIXXb":"1","zChJod":"%.@.]","Yllh3e":"%.@.1641105444967722,178657810,1996762967]","w2btAe":"%.@.\"101199000968061242063\",\"101199000968061242063\",\"0\",null,null,null,1]","QrtxK":"0","eptZe":"/wizrpcui/_/WizRpcUi/","S06Grb":"101199000968061242063"};window.IJ_values=[false,true,true,true,false,false,false,24,"none",true,"1px 1px 15px 0px #171717",false,"rgba(255,255,255,.54)","rgba(255,255,255,.26)","#000","rgba(0,0,0,.3)",true,"none","#424548",true,false,false,true,false,"#609beb","#8ab4f8","1px 1px 15px 0px #171717",true,false,36,24,28,6,true,false,false,false,false,false,"#3c4043",10,true,false,false,"#202124","#e8eaed",false,"#303134","0px 5px 26px 0px rgba(0,0,0,0.5), 0px 20px 28px 0px rgba(0,0,0,0.5)","#4285f4",false,true,false,"#8ab4f8",false,true,false,false,"#fff","#4487f6","#48a1ff","#a4c2ff","#219540","#41b85f","#4eb66e","#ff7d70","#ff897e","#000","#1aa863","#212327","#050505","#0a0a0a","#111","#9aa0a6","#8a8a8a","#bdc1c6","#bdc1c6","#bdc1c6","#000","#8a4a00","#824300","#b85100","18px","#3c4043","#e8eaed","#e8eaed","#3c4043",14,"#e8eaed",40,"#e8eaed",false,"#8a8a8a","#ddd","#ff7a8e","#bdc1c6","arial,sans-serif-medium,sans-serif","arial,sans-serif","#bdc1c6","#292e36","#bdc1c6","#868b90","#48a1ff",false,false,false,false,false,false,true,false,false,false,"1px 1px 15px 0px #171717",false,false,"#3c4043","rgba(255,255,255,.26)","#9aa0a6","#bdc1c6","rgba(204,204,204,.15)","rgba(204,204,204,.25)","rgba(102,102,102,.2)","rgba(102,102,102,.4)","rgba(255,255,255,.12)","#3c4043","#fff","rgba(0,0,0,.3)","#000","#bdc1c6","#000","Roboto,RobotoDraft,Helvetica,Arial,sans-serif","14px","500","500","pointer","0 1px 1px rgba(0,0,0,.16)",true,24,"#000","1px 1px 15px 0px #171717","#dadce0",200,true,true,false,false,true,true,false,true,14,"#202124","#303134",false,"1px solid #3c4043","none","arial,sans-serif-medium,sans-serif","Google Sans,arial,sans-serif-medium,sans-serif","#3c4043","1px solid #3c4043","1px solid #5f6368","rgba(255,255,255,0.1)","#3c4043","#202124","#8ab4f8","#3c4043","#bdc1c6","#9aa0a6",false,true,true,false,false,false,false,false,false,false,false,false,true,false,false,false,true,false,false,false,false,false,false,false,"8px","#3c4043",false,true,false,"%.@.\"101199000968061242063\",\"101199000968061242063\",\"0\",null,null,null,1]","0","%.@.null,1,1,null,[null,757,1440]]","LdSGxtnn+iuvNIwJ+JnlLg\u003d\u003d","%.@.\"#424548\"]","%.@.0]","%.@.0]","%.@.\"0px 5px 26px 0px rgba(0,0,0,0.5),0px 20px 28px 0px rgba(0,0,0,0.5)\",\"#303134\"]","%.@.0,null,null,36,28,6,0.3,null,14,null,null,null,null,null,\"#bdc1c6\",\"#9aa0a6\",null,\"#bdc1c6\",null,null,null,null,null,null,\"#1a73e8\",\"#fabb05\",\"#fff\",\"#1a73e8\",\"#d1d1d1\",\"#fff\",null,null,null,14,500,\"#51a6ff\",null,\"#8ab4f8\",\"#303134\"]",null,"%.@.[],0,null,1,1]","zh-Hans-TW","%.@.\"13px\",\"16px\",\"11px\",13,16,11,\"8px\",8,20]","zh_Hans_TW","%.@.\"10px\",10,\"16px\",16,\"18px\"]","%.@.\"14px\",14]","%.@.40,32,14]",null,"%.@.\"1px 1px 15px 0px #171717\"]","%.@.0,\"14px\",\"500\",\"500\",\"0 1px 1px rgba(0,0,0,.16)\",\"pointer\",\"#fff\",\"rgba(255,255,255,.26)\",\"#9aa0a6\",\"#bdc1c6\",\"rgba(204,204,204,.15)\",\"rgba(204,204,204,.25)\",\"rgba(102,102,102,.2)\",\"rgba(102,102,102,.4)\",\"#1aa863\",\"#4487f6\",\"#a4c2ff\",\"#ff7d70\",\"#8a4a00\",\"#111\",\"#050505\",\"#bdc1c6\",\"#4f861f\",\"rgba(255,255,255,.12)\",null,\"#000\",\"rgba(0,0,0,.3)\",\"#000\",\"#bdc1c6\",\"#000\",null,0]","%.@.\"20px\",\"500\",\"400\",\"13px\",\"15px\",\"15px\",\"Roboto,RobotoDraft,Helvetica,Arial,sans-serif\",\"24px\",\"400\",\"32px\",\"24px\"]",false,"","%.@.null,null,null,null,\"20px\",\"20px\",\"18px\",\"40px\",\"36px\",\"32px\",null,null,null,null,null,null,\"#202124\",null,null,null,\"#202124\",null,null,null,\"rgba(138,180,248,0.24)\",null,\"rgba(138,180,248,0.24)\",null,null,\"16px\",\"12px\",\"8px\",\"4px\",\"#202124\",\"rgba(138,180,248,0.24)\",\"#d2e3fc\",\"transparent\",\"#8ab4f8\",\"#5f6368\",\"999rem\",\"8px\",\"#d2e3fc\",\"transparent\",\"#dadce0\",\"#5f6368\",\"#d2e3fc\",\"transparent\",\"#8ab4f8\",\"#5f6368\",\"999rem\",\"Google Sans,arial,sans-serif-medium,sans-serif\",\"20px\",\"14px\",\"500\"]","%.@.\"#bdc1c6\",\"#bdc1c6\",\"#8ab4f8\",null,\"#9aa0a6\",\"#8ab4f8\",\"#c58af9\",null,null,\"#202124\",\"#8ab4f8\",\"#202124\",\"#394457\",\"#d2e3fc\",\"#303134\",\"#bdc1c6\",\"#fff\",\"#3c4043\",\"#202124\",\"#fff\",\"#202124\",\"#fff\",\"#81c995\",\"#f28b82\",\"#fdd663\",\"#3c4043\",\"#202124\",\"rgba(0,0,0,0.6)\",\"#bdc1c6\",\"#3c4043\"]","%.@.null,\"none\",null,\"0px 1px 3px hsla(0,0%,9%,0.24)\",null,\"0px 2px 6px hsla(0,0%,9%,0.32)\",null,\"0px 4px 12px hsla(0,0%,9%,0.9)\",null,null,\"1px solid  #5f6368\",\"none\",\"none\",\"none\"]","%.@.\"Google Sans,arial,sans-serif\",\"Google Sans,arial,sans-serif-medium,sans-serif\",\"arial,sans-serif\",\"arial,sans-serif-medium,sans-serif\",\"arial,sans-serif-light,sans-serif\"]","%.@.\"16px\",\"12px\",\"0px\",\"8px\",\"4px\",\"2px\",\"20px\",\"24px\"]","%.@.\"#8ab4f8\",\"#8ab4f8\"]","%.@.null,null,null,null,null,null,null,\"12px\",\"8px\",\"4px\",\"16px\",\"2px\",\"999rem\",\"0px\"]","%.@.\"700\",\"400\",\"underline\",\"none\",\"capitalize\",\"none\",\"uppercase\",\"none\",\"500\",\"lowercase\",\"italic\",\"-1px\",\"0.3px\"]","%.@.\"20px\",\"26px\",\"400\",\"Google Sans,arial,sans-serif\",null,\"arial,sans-serif\",\"14px\",\"400\",\"22px\",null,\"16px\",\"24px\",\"400\",\"Google Sans,arial,sans-serif\",null,\"Google Sans,arial,sans-serif\",\"60px\",\"48px\",\"-1px\",null,\"400\",\"Google Sans,arial,sans-serif\",\"36px\",\"400\",\"48px\",null,\"Google Sans,arial,sans-serif\",\"36px\",\"28px\",null,\"400\",null,\"arial,sans-serif\",\"24px\",\"18px\",null,\"400\",\"arial,sans-serif\",\"16px\",\"12px\",null,\"400\",\"arial,sans-serif\",\"22px\",\"16px\",null,\"400\",\"arial,sans-serif\",\"26px\",\"20px\",null,\"400\",\"arial,sans-serif\",\"20px\",\"16px\",null,\"400\",\"arial,sans-serif\",\"18px\",\"14px\",null,\"400\",\"Google Sans,arial,sans-serif\",\"32px\",\"24px\",null,\"500\"]","%.@.4]","%.@.\"14px\",14,\"16px\",16,\"0\",0,\"none\",632,\"1px solid #3c4043\",\"normal\",\"normal\",\"#9aa0a6\",\"12px\",\"1.34\",\"1px solid #3c4043\",\"none\",\"0\",\"none\",\"none\",\"none\",\"none\",\"6px\"]","%.@.\"0\"]","%.@.\"rgba(0,0,0,0.0)\",\"rgba(0,0,0,0.54)\",\"rgba(0,0,0,0.8)\",\"rgba(248, 249, 250, 0.85)\",\"#202124\",\"#dadce0\",\"rgba(218, 220, 224, 0.0)\",\"rgba(218, 220, 224, 0.7)\",\"#dadce0\",\"#f8f9fa\",\"#000\",\"#1a73e8\",\"#dadce0\",\"#fff\",\"#fff\",\"#e8eaed\"]","%.@.\"#dddee1\",\"#868b90\",\"#bdc1c6\",\"#bcc0c3\",\"#000\",\"rgba(0,0,0,.7)\",28,24,26,20,16,-2,0,-4,2,0,0,24,20,20,14,12]","%.@.\"20px\",20,\"14px\",14,\"#e8eaed\"]","troy450409405@gmail.com",true,"101199000968061242063","%.@.1]"];})();(function(){google.llirm='400px';google.ldi={};google.pim={};})();
window.jsl=window.jsl||{};window.jsl.dh=function(d,e,c){try{var f=document.getElementById(d);if(f)f.innerHTML=e,c&&c();else{var a={id:d,script:String(!!c),milestone:String(google.jslm||0)};google.jsla&&(a.async=google.jsla);var g=document.createElement("div");g.innerHTML=e;var b=g.children[0];b&&(a.tag=b.tagName,a["class"]=String(b.className||null),a.name=String(b.getAttribute("jsname")));google.ml(Error("Missing ID."),!1,a)}}catch(h){google.ml(h,!0,{"jsl.dh":!0})}};(function(){var x=true;
google.jslm=x?2:1;})();google.x(null, function(){(function(){(function(){google.csct={};google.csct.ps='AOvVaw17ag9mz-2UL3tGGKniglcH\x26ust\x3d1641191845010685';})();})();(function(){(function(){google.csct.rw=true;})();})();(function(){(function(){google.csct.rl=true;})();})();(function(){google.drty&&google.drty(undefined,true);})();});google.drty&&google.drty(undefined,true);`,
		},
		{
			name: "13",
			args: args{
				content: "https://example.com:8080/demo",
				oldHost: "example.com:8080",
				newHost: "localhost",
			},
			want: "http://localhost/demo",
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
