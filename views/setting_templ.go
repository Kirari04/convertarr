// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"encoder/helper"
	"encoder/layouts"
	"encoder/t"
	"fmt"
	"html"
	"runtime"
	"time"
)

func Setting(Ctx t.TemplCtx, Title string, v t.SettingValidator) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<article class=\"message\"><div class=\"message-header\"><p>Setting</p></div><form method=\"POST\" action=\"/setting\" class=\"message-body\"><h2 class=\"subtitle is-4\">Security</h2><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableAuthentication\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableAuthentication) == "checked" || helper.PStrToStr(v.EnableAuthentication) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Authentication</label></div><div class=\"field\"><div class=\"control\"><label class=\"radio\"><input type=\"radio\" name=\"AuthenticationType\" value=\"form\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.AuthenticationType) == "form" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Form-Html</label> <label class=\"radio\"><input type=\"radio\" name=\"AuthenticationType\" value=\"basic\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.AuthenticationType) == "basic" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Basic-Auth</label></div></div><div class=\"field\"><a href=\"/setting/user\" class=\"button is-info\">User Settings</a></div><br><h2 class=\"subtitle is-4\">Folder Scanning</h2><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableAutomaticScanns\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableAutomaticScanns) == "checked" || helper.PStrToStr(v.EnableAutomaticScanns) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Automatic Scanns</label></div><div class=\"field\"><div class=\"control\"><span class=\"select\"><select name=\"AutomaticScannsInterval\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, min := range []int{2, 5, 15, 30, 60, 120, 720, 1440} {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.EscapeString(fmt.Sprint(min))))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if v.AutomaticScannsInterval == min {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" selected")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(templ.EscapeString(fmt.Sprintf("%s", time.Duration(min)*time.Minute)))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/setting.templ`, Line: 75, Col: 84}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></span> Automatic Scanns Interval</div></div><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"AutomaticScannsAtStartup\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.AutomaticScannsAtStartup) == "checked" || helper.PStrToStr(v.AutomaticScannsAtStartup) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Automatic Scanns At Startup</label></div><div class=\"field\"><a href=\"/setting/folder\" class=\"button is-info\">Folder Settings</a></div><h2 class=\"subtitle is-4\">File Encoding</h2><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableEncoding\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableEncoding) == "checked" || helper.PStrToStr(v.EnableEncoding) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Encoding</label></div><div class=\"field\"><div class=\"control\"><input class=\"input\" type=\"number\" name=\"EncodingThreads\" style=\"max-width: 400px;\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.EscapeString(fmt.Sprint(v.EncodingThreads))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" min=\"0\" max=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprint(runtime.NumCPU())))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> Encoding Threads (0 = use all, max = ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(runtime.NumCPU()))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/setting.templ`, Line: 120, Col: 73}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(") <br>When using Hevc codec, threads have a different meaning because of pools: <a href=\"https://trac.ffmpeg.org/ticket/3730\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(html.EscapeString("https://trac.ffmpeg.org/ticket/3730"))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/setting.templ`, Line: 124, Col: 65}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</a></div></div><div class=\"field\"><div class=\"control\"><input class=\"input\" type=\"number\" name=\"EncodingCrf\" style=\"max-width: 400px;\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.EscapeString(fmt.Sprint(v.EncodingCrf))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" min=\"1\" max=\"50\"> Encoding Crf (1-50)</div></div><div class=\"field\"><div class=\"control\"><input class=\"input\" type=\"number\" name=\"EncodingResolution\" style=\"max-width: 400px;\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.EscapeString(fmt.Sprint(v.EncodingResolution))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" min=\"100\" max=\"5000\"> Encoding Resolution (100-5000)</div></div><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableHevcEncoding\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableHevcEncoding) == "checked" || helper.PStrToStr(v.EnableHevcEncoding) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Hevc Encoding</label></div><div class=\"field\"><div class=\"control\"><input class=\"input\" type=\"number\" name=\"EncodingMaxRetry\" style=\"max-width: 400px;\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.EscapeString(fmt.Sprint(v.EncodingMaxRetry))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" min=\"0\" max=\"999\"> Encoding Max Retry (0-999)</div></div><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableNvidiaGpuEncoding\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableNvidiaGpuEncoding) == "checked" || helper.PStrToStr(v.EnableNvidiaGpuEncoding) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Nvidia Gpu Encoding</label></div><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableAmdGpuEncoding\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableAmdGpuEncoding) == "checked" || helper.PStrToStr(v.EnableAmdGpuEncoding) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Amd Gpu Encoding</label></div><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableImageComparison\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if helper.PStrToStr(v.EnableImageComparison) == "checked" || helper.PStrToStr(v.EnableImageComparison) == "on" {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" checked")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Enable Image Comparison</label></div><div class=\"field\"><div class=\"control\"><button type=\"submit\" class=\"button is-primary\">Save</button></div></div></form></article>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layouts.Default(Ctx, Title).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
