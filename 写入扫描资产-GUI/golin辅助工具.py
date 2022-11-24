import tkinter as tk
import tkinter.messagebox 
import os
from tkinter import ttk

root = tk.Tk()
root.title("Golin辅助-高业尚v2")
def set_win_center(root, curWidth='', curHight=''):
     if not curWidth:
         '''获取窗口宽度，默认200'''
         curWidth = root.winfo_width()
     if not curHight:
         '''获取窗口高度，默认200'''
         curHight = root.winfo_height()
     scn_w, scn_h = root.maxsize()
     cen_x = (scn_w - curWidth) / 2
     cen_y = (scn_h - curHight) / 2
     # print(cen_x, cen_y)
 
     size_xy = '%dx%d+%d+%d' % (curWidth, curHight, cen_x, cen_y)
     root.geometry(size_xy)

#root.resizable(False, False)  # 窗口不可调整大小
root.update()  # 必须
set_win_center(root, 280, 280)

xVariable = tkinter.StringVar()
# 设置标签信息
label1 = tk.Label(root, text='服务器名称')
label1.grid(row=0, column=0)
label2 = tk.Label(root, text='IP')
label2.grid(row=1, column=0)
label3 = tk.Label(root, text='用户名')
label3.grid(row=2, column=0)
label4 = tk.Label(root, text='密码')
label4.grid(row=3, column=0)
label5 = tk.Label(root, text='端口')
label5.grid(row=4, column=0)
label6 = tk.Label(root, text='写入文件')
label6.grid(row=5, column=0)
label7 = tk.Label(root, text='运行模式')
label7.grid(row=6, column=0)


# 创建输入框
entry1 = tk.Entry(root)
entry1.grid(row=0, column=1, padx=10, pady=5)
entry2 = tk.Entry(root)
entry2.grid(row=1, column=1, padx=10, pady=5)
entry3 = tk.Entry(root)
entry3.grid(row=2, column=1, padx=10, pady=5)
entry4 = tk.Entry(root,show='*')
entry4.grid(row=3, column=1, padx=10, pady=5)
entry5 = tk.Entry(root)
entry5.grid(row=4, column=1, padx=10, pady=5)
# entry6 = tk.Entry(root)
# entry6.grid(row=5, column=1, padx=10, pady=5)
# entry7 = tk.Entry(root)

entry7 = ttk.Combobox(root, textvariable=xVariable,width=8,height=8)     # #创建下拉菜单
entry7.grid(row=6, column=1, padx=10, pady=5)
entry7["value"] = ("linux", "mysql", "redis","postgresql","windows","oracle")    # #给下拉菜单设定值
entry7.current(0)    # #设定下拉菜单的默认值为第0个


# def xFunc(event):
#     print(entry7.get())  # #获取选中的值方法1
#     print(xVariable.get())  # #获取选中的值方法2
#
# entry7.bind("<<ComboboxSelected>>", xFunc)  # #给下拉菜单绑定事件
# 创建按键
def show():
    if "~" in entry1.get() or "~" in entry4.get():
        tk.messagebox.showerror('错误！', '字段中不能包含~符号！')
        return
    name = entry1.get()+"~"
    ip = entry2.get()+"~"
    user = entry3.get()+"~"
    passwd = entry4.get()+"~"
    port = entry5.get()
    # firename =  entry7.get()
    global firename
    if entry7.get()=="linux":
        firename="ip.txt"
    if entry7.get()=="mysql":
        firename="mysql.txt"
    if entry7.get()=="redis":
        firename="redis.txt"
    if entry7.get()=="postgresql":
        firename="postgresql.txt"
    # a = name+ip+user+passwd+port
    #拼接目录
    pwd = os.getcwd()
    pwffire = os.path.join(pwd,firename)
    #写入资产数据
    with open(pwffire, 'a+',encoding="utf-8") as f:
        size = os.path.getsize(pwffire)
        if size == 0:
            name = entry1.get() + "~"
            a = name + ip + user + passwd + port
        else:
            name = "\n" + entry1.get() + "~"
            a = name + ip + user + passwd + port
        f.write(a)
        print(a)
##开始运行
def run():
    pwd = os.getcwd()
    check = os.path.isfile(os.path.join(pwd, "golin.exe"))
    if check != True:
        tk.messagebox.showerror('错误！', '当前目录下不存在golin.exe程序，可通过https://github.com/selinuxG/Golin下载')
        return
    if entry7.get()=="oracle":
        tk.messagebox.showinfo('提示', '如果调用oracle模式，首先确保本机有oracle的客户端或者官方sdk\n需要在cmd下运行\ngolin.exe -run oracle  -con system/oracle@1.1.2.135:1521/sid -name oracle')
        return
    runtype = f"golin.exe -run {(entry7.get())}"
    pwffire = os.path.join(pwd,runtype)
    os.system(pwffire)
    successpath = os.path.join(pwd,"采集完成目录")
    if os.path.isdir(successpath):
        if entry7.get()!="windows":
            os.system("explorer "+successpath)
    tk.messagebox.showinfo('提示','采集完成💕')

#清空文件
def delfile():
    if entry7.get()=="linux":
        firename="ip.txt"
    if entry7.get()=="mysql":
        firename="mysql.txt"
    if entry7.get()=="redis":
        firename="redis.txt"
    if entry7.get()=="postgresql":
        firename="postgresql.txt"
    pwd = os.getcwd()
    pwffire = os.path.join(pwd, firename)
    try:
        os.remove(pwffire)
    except:
        pass


button1 = tk.Button(root, text='增加资产信息', background="#7FFFD4",command=show).grid(row=5, column=1,
                                             padx=30, pady=5)
button2 = tk.Button(root, text='退出程序', background="#FFC0CB",command=root.quit).grid(row=7, column=0,
                                          padx=30, pady=5)
button3 = tk.Button(root, text='运行采集功能', background="#7FFFD4",command=run).grid(row=7, column=1,
                                            padx=30, pady=5)
button4 = tk.Button(root, text='清空文件', background="#FFC0CB",command=delfile).grid(row=5, column=0,
                                            padx=30, pady=5)

tk.mainloop()

