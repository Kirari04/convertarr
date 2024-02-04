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
	"encoder/layouts"
	"encoder/t"
	"fmt"
)

func chartData(resources t.Resources, maxResourcesHistory int) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_chartData_31e2`,
		Function: `function __templ_chartData_31e2(resources, maxResourcesHistory){const optimizeChartSeries = (data) => {
		if(!data){
			return [0]
		}
		const maxLen = 24;
		const currentLen = data.length;
		if(currentLen <= maxLen) {
			return data;
		}
		const ratio = Math.ceil(currentLen / maxLen)
		let newData = []
		for(let i = 0; i < maxLen; i++){
			let dataSlice = [];
			for(let x = 0; x < ratio; x++){
				if(data[x * i]){
					dataSlice.push(data[x * i])
				}
			}
			const sum = dataSlice.reduce((a, b) => a + b, 0);
			const avg = (sum / dataSlice.length) || 0;
			newData.push(avg)
		}
		return newData;
	}

	var options = {
		chart: {
			type: "line",
			height: "300px",

		},
		series: [
			{
				name: "Cpu Ussage",
				data: optimizeChartSeries(resources.Cpu),
			},
			{
				name: "Memory Ussage",
				data: optimizeChartSeries(resources.Mem),
			}
		],
		xaxis: {
			categories: Array.from(Array(maxResourcesHistory).keys()),
		},
		yaxis: {
            labels: {
                formatter: function (value) {
                    return ` + "`" + `${Math.round(value)} %` + "`" + `
                },
            },
            max: 100,
        },
        tooltip: {
            y: {
                formatter: function (value) {
                    return ` + "`" + `${Math.round(value)} %` + "`" + `
                },
            },
        }
	}
	var chart = new ApexCharts(document.querySelector("#chart"), options);

	chart.render();

	var options2 = {
		chart: {
			type: "line",
			height: "300px",

		},
		series: [
			{
				name: "NetOut",
				data: optimizeChartSeries(resources.NetOut),
			},
			{
				name: "NetIn",
				data: optimizeChartSeries(resources.NetIn),
			}
		],
		xaxis: {
			categories: Array.from(Array(maxResourcesHistory).keys()),
		},
		yaxis: {
            labels: {
                formatter: function (value) {
                    return ` + "`" + `${humanFileSize(value)}/s` + "`" + `
                },
            },
        },
        tooltip: {
            y: {
                formatter: function (value) {
                    return ` + "`" + `${humanFileSize(value)}/s` + "`" + `
                },
            },
        }
	}
	var chart2 = new ApexCharts(document.querySelector("#chart2"), options2);

	chart2.render();

	function humanFileSize(bytes, si = false, dp = 1) {
		const thresh = si ? 1000 : 1024;

		if (Math.abs(bytes) < thresh) {
			return bytes + ' B';
		}

		const units = si
			? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
			: ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
		let u = -1;
		const r = 10 ** dp;

		do {
			bytes /= thresh;
			++u;
		} while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


		return bytes.toFixed(dp) + ' ' + units[u];
	}
}`,
		Call:       templ.SafeScript(`__templ_chartData_31e2`, resources, maxResourcesHistory),
		CallInline: templ.SafeScriptInline(`__templ_chartData_31e2`, resources, maxResourcesHistory),
	}
}

func Index(Ctx t.TemplCtx, Title string, resources t.Resources, savedStorage string, encodedFiles string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<article class=\"message\"><div class=\"message-header\"><p>Server Resource Ussage</p></div><div class=\"message-body\"><nav class=\"level\"><div class=\"level-item has-text-centered\"><div><p class=\"heading\">Files in queue</p><p class=\"title\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(len(app.FilesToEncode)))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/index.templ`, Line: 146, Col: 60}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div><div class=\"level-item has-text-centered\"><div><p class=\"heading\">Saved storage</p><p class=\"title\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(savedStorage)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/index.templ`, Line: 152, Col: 38}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div><div class=\"level-item has-text-centered\"><div><p class=\"heading\">Encoded Files</p><p class=\"title\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(encodedFiles)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/index.templ`, Line: 158, Col: 38}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div></div></nav><div id=\"chart\"></div><div id=\"chart2\"></div></div></article>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = chartData(resources, app.MaxResourcesHistory).Render(ctx, templ_7745c5c3_Buffer)
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
