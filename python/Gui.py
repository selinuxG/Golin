# -*- coding: utf-8 -*-
import os
import subprocess
import tkinter as tk
import tkinter.messagebox
import webbrowser as web
from subprocess import run
from tkinter import ttk, filedialog
import requests

cmdpath = ""  # 运行linux模式下指定的cmd位置
root = tk.Tk()
root.title("Golin辅助-高业尚v2")

# 下载golin
def wget():
    try:
        r = requests.get("https://github.com/selinuxG/Golin-cli/releases/download/v1.3/golin.exe")
        f = open("golin.exe", "wb")
        f.write(r.content)
        f.close()
        tk.messagebox.showinfo('提示', '下载完成💕')
    except:
        tk.messagebox.showerror('错误！', '下载失败!\n手动下载地址:https://github.com/selinuxG/Golin-cli/releases')

# 打开帮助
def help():
    url = 'https://github.com/selinuxG/Golin-cli'
    web.open(url)

# 读取自定义命令路径
def cmd():
    global cmdpath
    cmdpath = filedialog.askopenfilename()

# 重置输入框内容
def delval():
    global cmdpath
    entry1.delete(0, tkinter.END)
    entry2.delete(0, tkinter.END)
    entry3.delete(0, tkinter.END)
    entry4.delete(0, tkinter.END)
    entry5.delete(0, tkinter.END)
    cmdpath = ""
    entry7.current(0)

# cli模式
def cli():
    # 接收执行的命令并执行
    def cli_input():
        command = cmdrun.get()
        #先删除
        output_box.delete(1.0, tk.END)
        process = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True, text=True,
                                   encoding="utf-8", errors='replace')
        output, error = process.communicate()
        process.wait()
        if process.returncode != 0:
            print(error)
            output_box.insert(tk.END, "Error: " + error + "\n")
        else:
            output_box.insert(tk.END, output.strip() + "\n")

    root = tk.Toplevel()
    root.title("Cli模式")
    window_width,window_height = 855,180    # 设置新窗口的尺寸
    screen_width,screen_height = root.winfo_screenwidth(),root.winfo_screenheight()    # 获取屏幕的宽度和高度
    x_position,y_position = (screen_width // 2) - (window_width // 2),(screen_height // 2) - (window_height // 2)    # 计算新窗口的 x 和 y 坐标以使其居中
    root.geometry(f"{window_width}x{window_height}+{x_position}+{y_position}")
    # 输入命令提示
    tk.Label(root, text="输入命令", background="#40E0D0").grid(row=0, column=0, padx=10, pady=5)
    # 输入命令框
    cmdrun = tk.Entry(root,width=105,textvariable=tk.StringVar(value='golin '))
    cmdrun.grid(row=0, column=1)
    # 输入框提交
    tk.Button(root, text='run', background="#7FFFD4", command=cli_input).grid(row=0, column=2,padx=10,pady=5)
    # 输入命令提示
    tk.Label(root, text="结果展示", background="#40E0D0").grid(row=1, column=0, padx=10, pady=5)
    # 结果展示
    output_box = tk.Text(root, wrap=tk.WORD, height=10, width=105)
    output_box.grid(row=1, column=1)


menu1 = tk.Menu(root)
menu1.add_command(label="下载Golin", command=wget)
menu1.add_command(label="自定义命令", command=cmd)
menu1.add_command(label="重置", command=delval)
menu1.add_command(label="cli", command=cli)
menu1.add_command(label="help", command=help)

root.config(menu=menu1)
# 设置尺寸
sw = root.winfo_screenwidth()
sh = root.winfo_screenheight()
ww = 300
wh = 280
x = (sw - ww) / 2
y = (sh - wh) / 2
root.geometry("%dx%d+%d+%d" % (ww, wh, x, y))
root.update()  # 必须
# 下拉菜单的值
xVariable = tkinter.StringVar()
# 设置标签信息
label1 = tk.Label(root, text='定义名称', background="#40E0D0")
label1.grid(row=0, column=0)
label2 = tk.Label(root, text='连接地址', background="#40E0D0")
label2.grid(row=1, column=0)
label3 = tk.Label(root, text='连接用户', background="#40E0D0")
label3.grid(row=2, column=0)
label4 = tk.Label(root, text='连接密码', background="#40E0D0")
label4.grid(row=3, column=0)
label5 = tk.Label(root, text='连接端口', background="#40E0D0")
label5.grid(row=4, column=0)
label7 = tk.Label(root, text='运行模式')
label7.grid(row=6, column=0)

# 创建输入框
entry1 = tk.Entry(root)
entry1.grid(row=0, column=1, padx=10, pady=5)
entry2 = tk.Entry(root)
entry2.grid(row=1, column=1, padx=10, pady=5)
entry3 = tk.Entry(root)
entry3.grid(row=2, column=1, padx=10, pady=5)
entry4 = tk.Entry(root, show='*')
entry4.grid(row=3, column=1, padx=10, pady=5)
entry5 = tk.Entry(root)
entry5.grid(row=4, column=1, padx=10, pady=5)
entry7 = ttk.Combobox(root, textvariable=xVariable, width=8, height=8)  # #创建下拉菜单
entry7.grid(row=6, column=1, padx=10, pady=5)
entry7["value"] = ("linux", "mysql", "redis", "route")  # #给下拉菜单设定值
entry7.current(0)  # #设定下拉菜单的默认值为第0个

# 增加资产信息
def show():
    if "~" in entry1.get() or "~" in entry4.get():
        tk.messagebox.showerror('错误！', '字段中不能包含~符号！')
        return
    name = entry1.get() + "~"
    ip = entry2.get() + "~"
    user = entry3.get() + "~"
    passwd = entry4.get() + "~"
    port = entry5.get()
    # firename =  entry7.get()
    global firename
    if entry7.get() == "linux":
        firename = "linux.txt"
    if entry7.get() == "mysql":
        firename = "mysql.txt"
    if entry7.get() == "redis":
        firename = "redis.txt"
    if entry7.get() == "route":
        firename = "route.txt"

    pwd = os.getcwd()
    pwffire = os.path.join(pwd, firename)
    # 写入资产数据
    with open(pwffire, 'a+', encoding="utf-8") as f:
        size = os.path.getsize(pwffire)
        if size == 0:
            name = entry1.get() + "~"
            a = name + ip + user + passwd + port
        else:
            name = "\n" + entry1.get() + "~"
            a = name + ip + user + passwd + port
        f.write(a)
        print(a)


# 开始运行
def rungolin():
    global cmdpath
    pwd = os.getcwd()
    successpath = os.path.join(pwd, "采集完成目录")  # 采集完成目录
    runtype = f"golin.exe {entry7.get()}"  # 运行模式
    pwffire = os.path.join(pwd, runtype)  # 拼接golin+运行模式路径
    check = os.path.isfile(os.path.join(pwd, "golin.exe"))
    if not check:
        tk.messagebox.showerror('错误！', '当前目录下不存在golin.exe程序，可通过https://github.com/selinuxG/Golin下载')
        return

    # 运行linux模式下的自定义cmd命令
    if entry7.get() == "linux" and len(cmdpath) != 0:
        runtype = runtype + f" --cmd {cmdpath}"
        print(runtype)
        pwffire = os.path.join(pwd, runtype)  # 拼接golin+运行模式路径
        run(pwffire, shell=True)
        tk.messagebox.showinfo('提示', '自定义采集完成✔')
        if os.path.isdir(successpath):
            run("explorer " + successpath, shell=True)
        return
    # 调用其他模式
    # os.system("start "+pwffire)
    run(pwffire, shell=True)
    if os.path.isdir(successpath):
        if entry7.get() != "windows":
            run("explorer " + successpath, shell=True)
    tk.messagebox.showinfo('提示', '采集完成💕')


# 清空文件
def delfile():
    if entry7.get() == "linux":
        firename = "linux.txt"
    if entry7.get() == "mysql":
        firename = "mysql.txt"
    if entry7.get() == "redis":
        firename = "redis.txt"
    if entry7.get() == "route":
        firename = "route.txt"
    pwd = os.getcwd()
    pwffire = os.path.join(pwd, firename)
    try:
        os.remove(pwffire)
        tk.messagebox.showinfo("提示", f"{pwffire},清空完成!")
    except Exception as e:
        tk.messagebox.showwarning("警告", e)
        pass


tk.Button(root, text='增加资产信息', background="#7FFFD4", command=show).grid(row=5, column=1, padx=30,
                                                                              pady=5)
tk.Button(root, text='退出程序', background="#FFC0CB", command=root.quit).grid(row=7, column=0, padx=30,
                                                                               pady=5)
tk.Button(root, text='运行采集功能', background="#7FFFD4", command=rungolin).grid(row=7, column=1, padx=30,
                                                                                  pady=5)
tk.Button(root, text='清空文件', background="#FFC0CB", command=delfile).grid(row=5, column=0, padx=30, pady=5)
root.attributes("-toolwindow", 0)
tk.mainloop()
