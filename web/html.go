package web

// IndexHtml 单主机页面
func IndexHtml() string {
	return `
<!DOCTYPE html>
<html lang="zh">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Golin自动化平台_单主机(高业尚:版本)</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: Arial, sans-serif;
            background: linear-gradient(to right, #5886e2, #d55de0);
        }

        .container {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 90vh;
            /* 将min-height从100vh更改为90vh */
            margin-bottom: 20px;
            /* 增加底部外边距 */
        }

        /* 添加抖动动画 keyframes */
        @keyframes shake {

            0%,
            100% {
                transform: translateX(0);
            }

            10%,
            30%,
            50%,
            70%,
            90% {
                transform: translateX(-5px);
            }

            20%,
            40%,
            60%,
            80% {
                transform: translateX(5px);
            }
        }

        form {
            display: flex;
            flex-direction: column;
            gap: 30px;
            padding: 25px;
            border-radius: 12px;
            background-color: rgba(251, 250, 250, 0.9);
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24);
        }

        h1 {
            text-align: center;
            color: #5272f496;
            font-weight: bold;
            margin-bottom: 15px;
            font-size: 30px;
        }

        .input-row {
            display: flex;
            gap: 10px;
            /* 更改gap为10像素 */
            justify-content: center;
            /* 添加此项以居中输入框 */
        }

        input[type="text"],
        input[type="password"],
        input[type="number"] {
            width: 120px;
            /* 将输入框宽度减少到130像素 */
            padding: 10px;
            border: 1px solid #ccc;
            border-radius: 8px;
            font-size: 16px;
            transition: all 1s ease-in-out;
            animation: shake 30s linear;
            /* 修改抖动动画时间 */

        }

        input[type="text"]:focus,
        input[type="password"]:focus,
        input[type="number"]:focus {
            outline: none;
            border-color: #00ff88;
            box-shadow: 0 0 5px rgba(12, 217, 152, 0.775);
        }

        .select-container {
            position: relative;
        }

        select {
            width: 100%;
            /* 让选择栏宽度与整体输入框宽度一致 */
            padding: 10px;
            border-radius: 4px;
            border: 1px solid #ccc;
            font-size: 14px;
            appearance: none;
            background-color: #f8f9fa;
            /* 添加背景色 */
            color: #2d2d31;
            /* 更改文字颜色 */
            cursor: pointer;
            text-align-last: center;
            /* 使文本内容居中 */
            animation: shake 1s linear;
            /* 修改抖动动画时间 */

        }

        button {
            padding: 12px 25px;
            border: none;
            background-color: #8d61c7;
            color: #fff;
            font-size: 16px;
            font-weight: bold;
            border-radius: 4px;
            cursor: pointer;
            transition: all 0.3s ease-in-out;
            animation: shake 0.4s linear;
            /* 修改抖动动画时间 */

        }

        button:hover {
            background-color: #864ab1;
        }

        footer {
            position: fixed;
            bottom: 1rem;
            left: 50%;
            transform: translateX(-50%);
            font-size: 14px;
            text-align: center;
            color: #fff;
            /* 文字颜色更改为白色 */
        }

        footer a {
            text-decoration: none;
            color: rgb(82, 196, 54);
            transition: all 0.3s;
        }

        footer a:hover {
            color: rgba(82, 196, 54, 0.6);
            /* 鼠标悬停时降低透明度 */
        }

        .margin-top {
            margin-top: 20px;
            /* 自定义边距，可以根据需要调整 */
        }

        /* 此处以下为填写帮助 */
        .form-container {
            position: relative;
        }

        .help-btn {
            position: absolute;
            top: 1px;
            right: 1px;
            cursor: pointer;
            font-weight: bold;
            padding: 5px 12px;
            background-color: rgb(251, 250, 250);
            color: rgb(54, 222, 88);
            border-radius: 4px;
            transition: transform 0.3s ease-in-out;
            z-index: 1000;
            border: none;
        }

        .help-btn:hover {
            transform: scale(1.05);
        }

        .modal {
            display: none;
            position: fixed;
            z-index: 2000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            overflow: auto;
            background-color: rgba(0, 0, 0, 0.6);
        }

        .modal-content {
            background-color: #94f5ce;
            margin: 10% auto;
            padding: 20px;
            border: 1px solid #888;
            width: 500px;
            /* 设置宽度 */
            text-align: center;
            position: relative;
            max-height: 70%;
            /* 设置最大高度为视口的70% */
            overflow-y: auto;
            /* 添加滚动条以适应内容 */
        }


        .close {
            color: #aaa;
            position: absolute;
            top: 0;
            right: 14px;
            font-size: 28px;
            font-weight: bold;
            cursor: pointer;
        }

        .close:hover,
        .close:focus {
            color: black;
            text-decoration: none;
        }
    </style>
    <script>
        // 此处js代码只服务于填写手册
        document.addEventListener("DOMContentLoaded", function () {
            document.getElementById("helpBtn").addEventListener("click", function () {
                document.getElementById("myModal").style.display = "block";
            });

            document.getElementsByClassName("close")[0].addEventListener("click", function () {
                document.getElementById("myModal").style.display = "none";
            });

            window.addEventListener("click", function (event) {
                if (event.target == document.getElementById("myModal")) {
                    document.getElementById("myModal").style.display = "none";
                }
            });
        });

    </script>

</head>

<body>
    <div class="container">
        <form action="/golin/submit" method="post">
            <div class="form-container">
                <button type="button" id="helpBtn" class="help-btn">填写帮助</button>
            </div>

            <h1>单主机采集模式</h1>
            <div class="input-row">
                <input type="text" placeholder="名称" name="name" required>
                <input type="text" placeholder="IP" name="ip" required>
                <input type="text" placeholder="用户" name="user" list="user" required>
                <datalist id="user">
                    <option value="root">SSH</option>
                    <option value="root">MySQL</option>
                    <option value="null">Redis</option>
                    <option value="postgres">Pgsql</option>
                    <option value="sa">SqlServer</option>
                    <option value="admin">route</option>
                </datalist>
                <input type="password" placeholder="密码" name="password" required>
                <input type="number" placeholder="端口" name="port" id="port-input" list="port-options" min="1"
                    max="65535" required>
                <datalist id="port-options">
                    <option value="22">SSH</option>
                    <option value="3306">MySQL</option>
                    <option value="6379">Redis</option>
                    <option value="5432">Pgsql</option>
                    <option value="1433">SqlServer</option>
                    <option value="1521">Oracle</option>
                </datalist>
            </div>
            <div class="select-container">
                <select name="run_mode" required>
                    <option value="Linux">Linux</option>
                    <option value="MySQL">MySQL</option>
                    <option value="Redis">Redis</option>
                    <option value="pgsql">PostgreSQl</option>
                    <option value="sqlserver">SQLServer</option>
                    <option value="oracle">Oracle</option>
                    <option value="h3c">H3C</option>
                    <option value="huawei">Huawei</option>
                </select>
                <select name="down" class="margin-top" required>
                    <option value="down">下载</option>
                    <option value="preview">预览</option>
                </select>
            </div>
            <button type="submit">提交</button>
        </form>
    </div>

    <footer>
        version:版本 如觉得对自己有帮助点个星星吧~
        <a style="text-decoration: none;color: rgb(82, 196, 54);" href="https://github.com/selinuxG/Golin-cli"
            target="_blank">GitHub</a>
    </footer>
    <!-- 帮助手册 -->
    <div id="myModal" class="modal">
        <div class="modal-content">
            <span class="close">&times;</span>
            <!-- <h2>自定义内容标题</h2> -->
            <p>Pgsql模式默认连接的数据库为postgres;</p>
            <p>如果Redis模式无用户填写为null,代表空用户;</p>
            <p>Oracle模式默认连接的oid为orcl,如需更改名称后增加oid=名称;</p>
            <p>虽然咱们有此功能但是实在不建议使用此工具连接网络设备！</p>
        </div>
    </div>

</body>

</html>
`
}

// IndexFilehtml 多主机模式页面

func IndexFilehtml() string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Golin自动化平台_多主机(高业尚:版本)</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: #f2f2f2;
        }
        .container {
            width: 500px;
            background-color: #ffffff;
            padding: 40px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
            border-radius: 10px;
        }
        h1 {
            text-align: center;
            margin-bottom: 30px;
        }
        button[type="submit"],
        a.download-btn,
        label.upload-btn,
        a.single-host-mode-btn {
            display: block;
            text-align: center;
            text-decoration: none;
            color: #ffffff;
            padding: 15px;
            border-radius: 5px;
            border: none;
            cursor: pointer;
            margin-bottom: 20px;
            width: 100%;
            box-sizing: border-box;
        }
        a.download-btn {
            background-color: #3498db;
        }
        input[type="file"] {
            display: none;
        }
        label.upload-btn {
            background-color: #34495e;
        }
        button[type="submit"] {
            background-color: #27ae60;
        }
        a.single-host-mode-btn {
            background-color: #9b59b6;
        }
        select.mode-select {
            display: block;
            width: 100%;
            padding: 15px;
            margin-bottom: 20px;
        }
        footer {
            position: fixed;
            bottom: 1rem;
            left: 50%;
            transform: translateX(-50%);
            font-size: 14px;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>多主机模式</h1>
        <a href="/golin/index" class="single-host-mode-btn">单主机模式</a>
        <a href="/golin/modefile" class="download-btn">下载模板文件</a>
        <form action="/golin/submitfile" method="POST" enctype="multipart/form-data">
            <label for="file-upload" class="upload-btn">上传文件路径</label>
            <input type="file" id="file-upload" name="uploaded-file" onchange="document.querySelector('.upload-btn').textContent = this.files[0].name">
            <select class="mode-select" name="mode">
                <option value="Linux">Linux</option>
                <option value="Mysql">MySQL</option>
                <option value="Redis">Redis</option>
                <option value="pgsql">PostgreSQl</option>
                <option value="sqlserver">SqlServer</option>
                <option value="oracle">Oracle</option>
                <option value="h3c">H3C</option>
                <option value="huawei">Huawei</option>
            </select>
            <button type="submit" class="download-btn">提交任务</button>
        </form>
    </div>
    <footer>
        version:版本 如觉得对自己有帮助点个星星吧~ 
        <a style="text-decoration: none;color: rgb(82, 196, 54);" href="https://github.com/selinuxG/Golin-cli" target="_blank">GitHub</a>
    </footer>
</body>
</html>
`
}

func ErrorHtml() string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>提示信息</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }

    body {
      font-family: 'Arial', sans-serif;
      background-color: #f9f9f9;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      overflow: hidden;
    }

    .container {
      position: relative;
      text-align: center;
    }

    h1 {
      font-size: 10em;
      color: rgb(70, 179, 251);
      text-shadow: 2px 4px 4px rgba(0, 0, 0, 0.15);
      animation: bounce 1s ease infinite;
    }

    p {
      font-size: 1.5em;
      margin-top: -20px;
      margin-bottom: 1em;
      color: rgb(102, 102, 102);
    }

    a {
      display: inline-block;
      background-color: #3b97d3;
      padding: 12px 24px;
      font-size: 1em;
      color: #fff;
      text-decoration: none;
      border-radius: 50px;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
      transition: background-color 0.3s ease;
    }

    a:hover {
      background-color: #3b87c3;
    }

    @keyframes bounce {
      0%, 20%, 60%, 100% {
        -webkit-transform: translateY(0);
                transform: translateY(0);
      }
      40% {
        -webkit-transform: translateY(-20px);
                transform: translateY(-20px);
      }
      80% {
        -webkit-transform: translateY(-10px);
                transform: translateY(-10px);
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>status</h1>
    <p>errbody</p>
    <a href="/golin/gys">返回首页</a>
  </div>
</body>
</html>
`
}

// GolinHomeHtml 首页
func GolinHomeHtml() string {
	return `
<!DOCTYPE html>
<!DOCTYPE html>
<html lang="zh-CN">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Golin Web</title>
    <style>
        body {
            font-family: Arial, "微软雅黑", sans-serif;
            margin: 0;
            padding: 0;
            background: linear-gradient(120deg, #e06645 0%, #007bff 100%);
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
        }

        .container {
            width: 100%;
            max-width: 960px;
            margin: 0 auto;
            padding: 20px;
            display: flex;
            flex-direction: column;
            align-items: center;
            background-color: rgba(255, 255, 255, 0.8);
            border-radius: 10px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }

        h1 {
            text-align: center;
            font-size: 3rem;
            color: #333;
            margin-bottom: 50px;
        }

        .btn-group {
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            gap: 20px;
        }

        .btn {
            display: inline-block;
            padding: 10px 20px;
            border-radius: 5px;
            font-size: 1.2rem;
            text-align: center;
            text-decoration: none;
            background-color: #4c5fe0;
            color: rgb(252, 252, 252);
            transition: background-color 0.3s, box-shadow 0.3s, transform 0.2s ease-out;
            box-shadow: 0 4px 14px 0 rgba(65, 135, 214, 0.39);
        }

        .btn:hover {
            background-color: #18e27d;
            box-shadow: 0 6px 20px rgba(100, 255, 180, 0.5);
            transform: scale(1.4);
        }

        .footer {
            text-align: center;
            margin-top: 50px;
        }

        .footer a {
            text-decoration: none;
            color: #2dcd57;
        }

        .footer a:hover {
            text-decoration: underline;
        }
    </style>
</head>

<body>
    <div class="container">
        <h1>Golin Web</h1>
        <div class="btn-group">
            <a href="/golin/indexfile" class="btn" target="_blank">多主机采集模式</a>
            <a href="/golin/index" class="btn" target="_blank">单主机采集模式</a>
            <a href="/golin/history" class="btn" target="_blank">历史记录</a>
            <a href="https://github.com/selinuxG/Golin-cli" target="_blank" class="btn">帮助手册</a>
            <a href="https://ihuace.yuque.com/org-wiki-ihuace-fyorr3/totgpt" target="_blank" class="btn">作业指导书</a>
            <a href="/golin/update" class="btn">检查更新</a>
        </div>
        <div class="footer">
            <p>version:版本 如觉得对自己有帮助点个星星吧~ <a href="https://github.com/selinuxG/Golin-cli" target="_blank">Github</a></p>
        </div>
    </div>
</body>

</html>
`
}

func GolinHistoryIndex() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>历史任务</title>
	<style>
	body {
		font-family: Arial, sans-serif;
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100vh;
		background-color: #f1f1f1;
		margin: 0;
	}
    .table-wrapper {
        width: 80%;
        border-radius: 10px;
        overflow: hidden;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
        transition: all 0.3s ease;
    }
    
    table {
        width: 100%;
        border-collapse: separate;
        background-color: white;
        margin: auto;
    }

    th,
    td {
        padding: 12px 15px;
        text-align: left;
        border-bottom: 1px solid #e0e0e0;
    }

    th {
        background-color: #3f51b5;
        color: white;
        font-weight: bold;
    }

    tr:nth-child(even) {
        background-color: #f8f8f8;
    }

    tr:hover {
        background-color: #e8f0ff;
    }

    tbody {
        display: block;
        max-height: 600px;
        overflow-y: auto;
    }

    thead,
    tbody tr {
        display: table;
        width: 100%;
        table-layout: fixed;
    }

    .table-wrapper:hover {
        box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
        transform: translateY(-5px);
    }

</style>
</head>

<body>
    <div class="table-wrapper">
    <table>
        <thead>
            <tr>
                <th>序号</th>
                <th>名称</th>
                <th>IP</th>
                <th>用户</th>
                <th>端口</th>
                <th>类型</th>
                <th>时间</th>
                <th>状态</th>
            </tr>
        </thead>
        <tbody>
		主机列表
        </tbody>
    </table>
</body>
</html>
`
}
