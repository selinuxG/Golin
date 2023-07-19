package run

func mysqlhtml() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MySQL安全策略核查</title>
	<link rel="icon" href="https://s1.ax1x.com/2023/07/19/pC7B5sx.jpg" sizes="16x16">
    <style>
        body {
            display: grid;
            grid-template-columns: 1fr 200px;
            gap: 10px;
            font-family: Arial, sans-serif;
            position: relative;
        }

        table {
    		border-collapse: collapse;
   		 	margin-bottom: 20px;
    		width: 100%;
    		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    		table-layout: fixed;
    		word-wrap: break-word;
        }

        th,
        td {
            border: 1px solid #ddd;
            padding: 15px;
            text-align: left;
        }

        th {
            background-color: #007BFF;
            color: white;
            font-weight: bold;
        }

        tr:nth-child(even) {
            background-color: #f9f9f9;
        }

        tr:hover {
            background-color: #e6f2ff;
        }

        .watermark {
            font-size: 36px;
            color: rgba(128, 128, 128, 0.2);
            position: absolute;
            z-index: -1;
            transform: rotate(-30deg);
        }

        #toc {
            position: fixed;
            top: 20px;
            right: 30px;
            padding-left: 10px;
            background-color: #f8f9fa;
            padding: 10px;
            border: 1px solid #dee2e6;
            border-radius: 5px;
            box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }

        #toc ul {
            list-style-type: none;
            padding: 0;
        }

        #toc a {
            text-decoration: none;
            color: #333;
            display: block;
        }

        #toc a:hover {
            color: #007BFF;
        }
		.watermark {
			font-size: 36px;
			color: rgba(128, 128, 128, 0.2);
			position: absolute;
			z-index: 1000;
			transform: rotate(-30deg);
    	}
        pre {
            background-color: #f8f9fa;
            border: 1px solid #dee2e6;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            font-family: "Courier New", Courier, monospace;
			white-space: pre-wrap;
			word-break: break-word;
        }
    </style>

<body>

    <div id="content">
        <center><h1>替换名称_MySQL安全策略核查</h1></center>

        <h2 id="version">版本信息</h2>
        <table>
            <thead>
                <tr>
                    <th>版本</th>
                    <th>连接ID</th>
                </tr>
            </thead>
            <tbody>
                版本详细信息
            </tbody>
        </table>

        <h2 id="userinfo">用户信息</h2>
        <table>
            <thead>
                <tr>
                    <th>用户</th>
                    <th>主机</th>
                    <th>密码信息</th>
                    <th>密码加密插件</th>
                    <th>加密连接类型</th>
                    <th>是否锁定</th>
                    <th>过期时间</th>
                    <th>是否过期</th>
                    <th>上次修改密码时间</th>
                    <th>账户类别</th>
                    <th>具体权限</th>
                </tr>
            </thead>
            <tbody>
                用户详细信息
            </tbody>
        </table>
       
        <h2 id="password">密码复杂度策略</h2>
        <table>
            <thead>
                <tr>
                    <th>是否允许密码与账户同名</th>
                    <th>密码长度要求</th>
                    <th>大写字符长度要求</th>
                    <th>数字字符长度要求</th>
                    <th>特殊字符长度要求</th>
                    <th>级别</th>
                </tr>
            </thead>
            <tbody>
                密码复杂度详细信息
            </tbody>
        </table>

        <h2 id="password-exp">密码过期时间</h2>
        <table>
            <thead>
                <tr>
                    <th>default_password_lifetime</th>
                    <th>天数</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>default_password_lifetime</td>
                    <td>密码过期时间结果</td>
                </tr>
            </tbody>
        </table>

		<h2 id="password-lock">失败锁定策略</h2>
		<table>
            <thead>
                <tr>
                    <th>连续失败次数</th>
                    <th>最大延迟(毫秒)</th>
                    <th>最小延迟(毫秒)</th>
                </tr>
            </thead>
            <tbody>
                密码过期时间详细信息
            </tbody>
        </table>

		<h2 id="timeout">超时策略(默认不符合)</h2>
		<table>
            <thead>
                <tr>
                    <th>登录时连接超时</th>
                    <th>登录后连接超时</th>
                </tr>
            </thead>
            <tbody>
                超时策略详细信息
            </tbody>
        </table>

		<h2 id="log">日志相关</h2>
		<table>
            <thead>
                <tr>
                    <th>错误存放路径</th>
                    <th>查询日志开启状态</th>
                    <th>查询日志路径</th>
                    <th>查询日志存放方式</th>
                </tr>
            </thead>
            <tbody>
                日志相关详细信息
            </tbody>
        </table>

		<h2 id="plug">插件信息</h2>
		<table>
            <thead>
                <tr>
                    <th>名称</th>
                    <th>状态</th>
                    <th>类型</th>
                    <th>插件库文件名</th>
                    <th>许可类型</th>
                </tr>
            </thead>
            <tbody>
                插件信息详细信息
            </tbody>
        </table>

		<h2 id="variables">系统变量</h2>
		<table>
            <thead>
                <tr>
                    <th>Key</th>
                    <th>Value</th>
                </tr>
            </thead>
            <tbody>
                系统变量详细信息
            </tbody>
        </table>


    </div>

    <div id="toc">
        <h3>目录</h3>
        <ul>
            <li><a href="#version">版本信息</a></li>
            <li><a href="#userinfo">用户信息</a></li>
            <li><a href="#password">密码复杂度信息</a></li>
            <li><a href="#password-exp">密码过期时间</a></li>
            <li><a href="#password-lock">失败锁定策略</a></li>
            <li><a href="#timeout">登录连接超时策略</a></li>
            <li><a href="#log">日志信息</a></li>
            <li><a href="#plug">插件信息</a></li>
            <li><a href="#variables">系统变量</a></li>

        </ul>
    </div>
	<div id="watermark"></div>
</body>

</html>
<script>
    const watermarkNum = 50 // 生成水印数量
    build()

    function build(){
        for(var i = 0; i < watermarkNum; i++){
            addWatermark(i);
        }
    }

    function addWatermark(i){
        var watermark = document.getElementById("watermark");
        const top = i
        const left = random();
        const  html = '<div class="watermark" style="top: '+(top/watermarkNum)*100+'%; left: '+left+'%;">高业尚-SelinuxG</div>'
        watermark.insertAdjacentHTML('afterend',html);
    }

    function random(){
       return Math.floor(Math.random() * 70) ;
    }
</script>
`
}
