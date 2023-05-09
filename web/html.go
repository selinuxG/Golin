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
    </style>
	<script>
		function handleSubmit(event) {
		  var inputs = document.querySelectorAll('input[type=text], input[type=password]');
		  for (var i = 0; i < inputs.length; i++) {
			if (inputs[i].value.trim() === '') {
			  alert("请确保所有输入框都填写了内容!");
			  event.preventDefault();
			  return;
			}
		  }}
	</script>
</head>

<body>
    <form action="/golin/submit" method="post" onsubmit="handleSubmit(event)">
        <label for="name">名称：</label>
        <input type="text" id="name" name="name">

        <label for="ip">IP：</label>
        <input type="text" id="ip" name="ip">

        <label for="user">用户：</label>
        <input type="text" id="user" name="user">

        <label for="password">密码：</label>
        <input type="password" id="password" name="password">

        <label for="port">端口：</label>
        <input type="text" id="port" name="port">

        <label for="run_mode">运行模式：</label>
        <select id="run_mode" name="run_mode" style="width: 100%;">
            <option value="Linux">Linux</option>
            <option value="Mysql">MySQL</option>
            <option value="Redis">Redis</option>
        </select>

        <input type="submit" value="提交">
    </form>
</body>

</html>
`
}
