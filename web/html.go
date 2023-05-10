package web

func IndexHtml() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Golin自动化平台</title>

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
        <div class="title">资产运行表</div>

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
                <option value="Mysql">MySQL</option>
                <option value="Redis">Redis</option>
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
</body>

</html>
`
}
