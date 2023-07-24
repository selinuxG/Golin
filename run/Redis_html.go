package run

func redishtml() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Redis安全策略核查</title>
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
		.permissions {
			width: 350px;
			white-space: nowrap;
			overflow-x: auto;
		}
    </style>

<body>

    <div id="content">
        <center><h1>替换名称_Redis安全策略核查</h1></center>

        <h2 id="info">基本信息</h2>
        <table>
            <thead>
                <tr>
                    <th>版本信息</th>
                    <th>监听网卡地址</th>
                    <th>配置文件中密码信息</th>
                </tr>
            </thead>
            <tbody>
                基本信息详细信息
            </tbody>
        </table>

        <h2 id="timout">超时及锁定信息</h2>
        <table>
            <thead>
                <tr>
                    <th>超时时间</th>
                    <th>失败锁定</th>
                </tr>
            </thead>
            <tbody>
                超时时间详细信息
            </tbody>
        </table>


        <h2 id="port">端口信息</h2>
        <table>
            <thead>
                <tr>
                    <th>非加密端口</th>
                    <th>加密端口</th>
                    <th>加密协议</th>
                </tr>
            </thead>
            <tbody>
                端口信息详细信息
            </tbody>
        </table>

        <h2 id="log">日志信息</h2>
        <table>
            <thead>
                <tr>
                    <th>日志存储位置</th>
                    <th>日志等级</th>
                    <th>Acl-log最大存储信息条数</th>
                </tr>
            </thead>
            <tbody>
                日志信息详细信息
            </tbody>
        </table>


        <h2 id="acluser">acluser</h2>
		<pre><code>acluser信息详细信息</code></pre>

        <h2 id="info">info信息</h2>
		<pre><code>info信息详细信息</code></pre>



    </div>

    <div id="toc">
        <h3>目录</h3>
        <ul>
            <li><a href="#info">基本信息</a></li>
            <li><a href="#timeout">超时及锁定信息</a></li>
            <li><a href="#port">端口信息</a></li>
            <li><a href="#log">日志信息</a></li>
            <li><a href="#acluser">acl用户信息</a></li>
            <li><a href="#info">info信息</a></li>
        </ul>
    </div>
	<div id="watermark"></div>
</body>

</html>
<script>
    const watermarkNum = 30 // 生成水印数量
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
