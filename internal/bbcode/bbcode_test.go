package bbcode

import "testing"

func TestRenderHtml(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "[b]hello[/b]", "<b>hello</b>"},
		{"code is raw without a default heading", "[code][b]hello[/b][/code]", `<div style="direction: ltr; margin: 5px; padding: 3px; border: 1px solid black; font-weight: normal; font-family: Monaco,'Courier New',monospace; background-color: rgb(242, 242, 242); overflow: scroll;">[b]hello[/b]</div>`},
		{"code with heading", "[code=Code:]hello[/code]", `<b>Code:</b><br><div style="direction: ltr; margin: 5px; padding: 3px; border: 1px solid black; font-weight: normal; font-family: Monaco,'Courier New',monospace; background-color: rgb(242, 242, 242); overflow: scroll;">hello</div>`},
		{"code heading is escaped", "[code=<script>]hello[/code]", `<b>&lt;script&gt;</b><br><div style="direction: ltr; margin: 5px; padding: 3px; border: 1px solid black; font-weight: normal; font-family: Monaco,'Courier New',monospace; background-color: rgb(242, 242, 242); overflow: scroll;">hello</div>`},
		{"url value", "[url=example.com]Example[/url]", `<a href="http://example.com" target="_blank">Example</a>`},
		{"url body", "[url]https://example.com[/url]", `<a href="https://example.com" target="_blank">https://example.com</a>`},
		{"url value with apostrophe", "[url=http://osu.titanic.sh/web/maps/009%20Sound%20System%20-%20Trinity%20(Cut%20Version)%20(Mahogany)%20%5bAdi's%20Nostalgia%5d.osu]Adi's Nostalgia[/url]", `<a href="http://osu.titanic.sh/web/maps/009%20Sound%20System%20-%20Trinity%20(Cut%20Version)%20(Mahogany)%20%5bAdi&#39;s%20Nostalgia%5d.osu" target="_blank">Adi&#39;s Nostalgia</a>`},
		{"email", "[email]test@example.com[/email]", `<a href="mailto:test@example.com">test@example.com</a>`},
		{"invalid email", "[email]invalid[/email]", "invalid"},
		{"spoiler", "[spoiler]spoiler content[/spoiler]", `<span style="background-color: black;">spoiler content</span>`},
		{"notice", "[notice]Notice[/notice]", `<div style="background: none repeat scroll 0% 0% rgb(249, 247, 254); border: 1px solid rgb(225, 223, 231); margin: 6px; padding: 5px;">Notice</div>`},
		{"box title with spaces", "[box=Mapping Tools:]body[/box]", `<div class="spoiler"><div class="spoiler-head" onclick="return toggleSpoiler(this);">Mapping Tools:</div><div class="spoiler-body">body</div></div>`},
		{"image", "[img]https://example.com/a.png[/img]", `<img src="https://example.com/a.png" loading="lazy">`},
		{"invalid image", "[img]not-a-url[/img]", ""},
		{"size clamp", "[size=900]large[/size]", `<span style="font-size:800%;">large</span>`},
		{"size named", "[size=small]small[/size]", `<span style="font-size:85%;">small</span>`},
		{"quote", "[quote=Alice]hello[/quote]", `<div class="quotetitle">Alice wrote:</div><div class="quotecontent">hello</div>`},
		{"quote strips bbcode", "[quote=Alice]hello [b]bold[/b] [url=https://example.com]link[/url][/quote]", `<div class="quotetitle">Alice wrote:</div><div class="quotecontent">hello bold link</div>`},
		{"spoilerbox trims newlines", "[spoilerbox]\n\nhello\nworld\n\n[/spoilerbox]", `<div class="spoiler"><div class="spoiler-head" onclick="return toggleSpoiler(this);">SPOILER</div><div class="spoiler-body">hello<br />world</div></div>`},
		{"list", "[list][*] one[*] two[/list]", "<ul><li> one</li><li> two</li></ul>"},
		{"unknown line", "[Header]\ntext", `<div class="beatmap-header">Header</div><br />text`},
		{"timecode", "01:23:456", `<a class="beatmap-timecode" href="osu://edit/01:23:456">01:23:456</a>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderHtml(tt.input)
			if got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestMediaProxy(t *testing.T) {
	renderer := New(Options{
		ImageProxyBaseUrl:  "https://proxy.example.com",
		ImageProxySecret:   "secret",
		ValidImageServices: []string{"trusted.example.com"},
	})

	source := "https://untrusted.example.com/a.png"
	got := renderer.RenderHtml("[img]" + source + "[/img]")
	want := `<img src="https://proxy.example.com` + signUrl(source, []byte("secret")) + `" loading="lazy">`
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestMediaProxyTrustedService(t *testing.T) {
	renderer := New(Options{
		ImageProxyBaseUrl:  "https://proxy.example.com",
		ImageProxySecret:   "secret",
		ValidImageServices: []string{"trusted.example.com"},
	})

	source := "https://trusted.example.com/a.png"
	got := renderer.RenderHtml("[img]" + source + "[/img]")
	want := `<img src="https://trusted.example.com/a.png" loading="lazy">`
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestRendererOptions(t *testing.T) {
	renderer := New(Options{BaseUrl: "https://example.com"})
	got := renderer.RenderHtml("[profile=1]Alice[/profile]")
	want := `<a href="https://example.com/u/1">Alice</a>`
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestStrip(t *testing.T) {
	got := Strip("[b]hello[/b]\n[i]world[/i]", false)
	want := "hello\nworld"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
