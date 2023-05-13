package web

func IndexHtml() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Golin自动化平台_单主机(高业尚:版本)</title>

    <style>
        html,
        body {
            width: 100%;
            height: 100%;
        }
        
        body {
            font-family: Arial, sans-serif;
            background-image: linear-gradient(to right, #56c6a8, #FED6E3);
            margin: 0;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        
        form {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
            /* max-width: 400px; */
            width: 400px;
        }
        
        label,
        input[type=text],
        input[type=password],
        select {
            display: block;
            margin-bottom: 15px;
        }
        
        input[type=text],
        input[type=password],
        select {
            padding: 10px;
            border: 2px solid #ccc;
            border-radius: 5px;
            outline: none;
            width: 100%;
            box-sizing: border-box;
            font-size: 14px;
        }
        
        input[type=text]:focus,
        input[type=password]:focus,
        select:focus {
            border-color: #FF6B81;
        }
        
        input[type=submit] {
            background-color: #FF6B81;
            color: white;
            padding: 10px 20px;
            margin-top: 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            width: 100%;
            transition: background-color 0.2s;
        }
        
        input[type=submit]:hover {
            background-color: #47ff50;
        }
        
        .from-group {
            display: flex;
            border-bottom: 1px solid #ccc;
            padding-top: 10px;
        }
        
        .from-group label {
            width: 120px;
            text-align: left;
            margin-right: 10px;
        }
        
        .title {
            color: #43ddcb;
            text-align: center;
            font-size: 29px;
            font-weight: bold;
            margin-bottom: 14px;
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
    <script>
        function handleSubmit(event) {
            var obj = document.getElementById("operate");
            var index = obj.selectedIndex;
            var value = obj.options[index].value;
            if (value == "preview") {
                document.getElementById("golinForm").target = "_blank";
                // event.target.action = "/golin/submit";
            } else {
                document.getElementById("golinForm").target = "_self";
            }
            var inputs = document.querySelectorAll('input[type=text], input[type=password]');
            for (var i = 0; i < inputs.length; i++) {
                if (inputs[i].value.trim() === '') {
                    alert("请确保所有输入框都填写了内容!");
                    event.preventDefault();
                    return;
                }
            }
        }
    </script>
</head>

<body>
    <form id="golinForm" action="/golin/submit" method="post" onsubmit="handleSubmit(event)">
        <div class="title">单主机模式:版本</div>

        <div class="from-group">
            <label for="name">名称：</label>
            <input type="text" id="name" name="name">
        </div>
        <div class="from-group">
            <label for="ip">IP：</label>
            <input type="text" id="ip" name="ip">
        </div>
        <div class="from-group">
            <label for="user">用户：</label>
            <input type="text" id="user" name="user">
        </div>
        <div class="from-group">
            <label for="password">密码：</label>
            <input type="password" id="password" name="password">
        </div>
        <div class="from-group">
            <label for="port">端口：</label>
            <input type="text" id="port" name="port">
        </div>
        <div class="from-group">
            <label for="run_mode">模式：</label>
            <select id="run_mode" name="run_mode" style="width: 100%;">
                <option value="Linux">Linux</option>
                <option value="MySQL">MySQL</option>
                <option value="Redis">Redis</option>
                <option value="Route">网络设备</option>
            </select>
        </div>

        <div class="from-group">
            <label for="operate">操作：</label>
            <select id="operate" name="down" style="width: 100%;">
                <option value="down">下载</option>
                <option value="preview">预览</option>
            </select>
        </div>
        <input type="submit" value="提交">
    </form>
    <footer>
        如有用点个星星吧~
        <a style="text-decoration: none;color: rgb(82, 196, 54);" href="https://github.com/selinuxG/Golin-cli" target="_blank">GitHub</a>
    </footer>
</body>

</html>
`
}

// 多主机文件

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
        <h1>多主机模式:版本</h1>
        <a href="/golin/index" class="single-host-mode-btn">单主机模式</a>
        <a href="/golin/modefile" class="download-btn">下载模板文件</a>
        <form action="/golin/submitfile" method="POST" enctype="multipart/form-data">
            <label for="file-upload" class="upload-btn">上传文件路径</label>
            <input type="file" id="file-upload" name="uploaded-file" onchange="document.querySelector('.upload-btn').textContent = this.files[0].name">
            <select class="mode-select" name="mode">
                <option value="Linux">Linux</option>
                <option value="Mysql">MySQL</option>
                <option value="Redis">Redis</option>
            </select>
            <button type="submit" class="download-btn">提交任务</button>
        </form>
    </div>
    <footer>
        如有用点个星星吧~ 
        <a style="text-decoration: none;color: rgb(82, 196, 54);" href="https://github.com/selinuxG/Golin-cli" target="_blank">GitHub</a>
    </footer>
</body>
</html>
`
}
