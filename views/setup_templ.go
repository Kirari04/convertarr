// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"encoder/app"
	"encoder/helper"
	"encoder/layouts"
	"encoder/t"
)

func Setup(Error error, Title string, v t.SetupValidator) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<article class=\"message is-info\"><div class=\"message-header\"><p>Initial Run</p></div><form method=\"POST\" action=\"/setup\" class=\"message-body\">Welcome to ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(app.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/setup.templ`, Line: 16, Col: 25}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<br>Configure the following initial settings to get up and running:<br><br><h2 class=\"subtitle is-4\">Security</h2><div class=\"field\"><label class=\"checkbox\"><input type=\"checkbox\" name=\"EnableAuthentication\"")
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("> Basic-Auth</label></div></div><div class=\"field\"><label class=\"label\">Username</label><div class=\"control\"><input class=\"input\" type=\"text\" autocomplete=\"username\" placeholder=\"Username\" name=\"Username\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.EscapeString(helper.PStrToStr(v.Username))))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" style=\"max-width: 400px;\"></div></div><div class=\"field\"><label class=\"label\">Password</label><div class=\"control\"><input class=\"input\" type=\"password\" autocomplete=\"new-password\" placeholder=\"Password\" name=\"Password\" style=\"max-width: 400px;\"></div></div><div class=\"field\"><div class=\"control\"><button type=\"submit\" class=\"button is-primary\">Save & Setup</button></div></div></form></article>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = layouts.Default(Error, Title).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
