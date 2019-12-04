// Copyright (c) 2013 Mathias Monnerville
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package cov

const (
	ProjectUrl = "https://github.com/matm/gocov-html"
	htmlHeader = `<html>
    <head>
        <title>Coverage Report</title>
        %s
    </head>
    <body>
        <div id="doctitle">Coverage Report</div>
    `

	htmlFooter = `
    </body>
</html>`

	// Default stylesheet
	defaultCSS = `
        <style type="text/css">
            body { background-color: #fff; }
            table {
                margin-left: 10px;
                border-collapse: collapse;
            }
            td { 
                background-color: #fff; 
                padding: 2px;
            }
            table.overview td {
                padding-right: 20px;
            }
            td.percent, td.linecount { text-align: right; }
            div.package, #totalcov { 
                color: #fff;
                background-color: #375eab; 
                font-size: 16px;
                font-weight: bold;
                padding: 10px;
                border-radius: 5px 5px 5px 5px;
            }
            div.package, #totalcov { 
                float: right; 
                right: 10px;
            }
            #totalcov { 
                top: 10px;
                position: relative;
                background-color: #fff;
                color: #000;
                border: 1px solid #375eab;
                clear: both;
            }
            #summaryWrapper {
                position: fixed;
                top: 10px;
                float: right;
                right: 10px;

            }
            span.packageTotal {
                float: right;
                color: #000;
            }
            #doctitle { 
                background-color: #fff; 
                font-size: 24px;
                margin-top: 20px;
                margin-left: 10px;
                color: #375eab;
                font-weight: bold;
            }
            #about {
                margin-left: 18px;
                font-size: 10px;
            }
            table tr:last-child td {
                font-weight: bold;
            }
            .functitle, .funcname { 
                text-align: center; 
                font-size: 20px; 
                font-weight: bold; 
                color: 375eab; 
            }
            .funcname {
                text-align: left;
                margin-top: 20px;
                margin-left: 10px;
                margin-bottom: 20px;
                padding: 2px 5px 5px;
                background-color: #e0ebf5;
            }
            table.listing {
                margin-left: 10px;
            }
            table.listing td {
                padding: 0px;
                font-size: 12px;
                background-color: #eee; 
                vertical-align: top;
                padding-left: 10px;
                border-bottom: 1px solid #fff;
            }
            table.listing td:first-child {
                text-align: right;
                font-weight: bold;
                vertical-align: center;
            }
            table.listing tr.miss td {
                background-color: #FFBBB8;
            }
            table.listing tr:last-child td {
                font-weight: normal;
                color: #000;
            }
            table.listing tr:last-child td:first-child {
                font-weight: bold;
            }
            .info {
                margin-left: 10px;
            }
            .info code {
            }
            pre { margin: 1px; }
            pre.cmd { 
                background-color: #e9e9e9;
                border-radius: 5px 5px 5px 5px;
                padding: 10px;
                margin: 20px;
                line-height: 18px;
                font-size; 14px;
            }
            a { 
                text-decoration: none; 
                color: #375eab;
            }
            a:hover { text-decoration: underline; }
            p { margin-left: 10px; }
        </style>
        `

	overview = `<p>This is a coverage report created after analysis of the <code>%s</code> package. It 
        has been generated with the following command:</p><pre class="cmd">gocov test %s | gocov-html</pre>        <p>Here are the stats. Please select a function name to view its implementation and see what's left for testing.</p>`
)
