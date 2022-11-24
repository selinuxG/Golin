import tkinter as tk
import tkinter.messagebox 
import os
root = tk.Tk()
root.title("增加资产信息-author:高业尚")
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
set_win_center(root, 380, 250)

# 设置标签信息
label1 = tk.Label(root, text='服务器名称:')
label1.grid(row=0, column=0)
label2 = tk.Label(root, text='IP:')
label2.grid(row=1, column=0)
label3 = tk.Label(root, text='用户:')
label3.grid(row=2, column=0)
label4 = tk.Label(root, text='密码:')
label4.grid(row=3, column=0)
label5 = tk.Label(root, text='端口:')
label5.grid(row=4, column=0)
label6 = tk.Label(root, text='写入文件:')
label6.grid(row=5, column=0)


# 创建输入框
entry1 = tk.Entry(root)
entry1.grid(row=0, column=1, padx=10, pady=5)
entry2 = tk.Entry(root)
entry2.grid(row=1, column=1, padx=10, pady=5)
entry3 = tk.Entry(root)
entry3.grid(row=2, column=1, padx=10, pady=5)
entry4 = tk.Entry(root,show="*")
entry4.grid(row=3, column=1, padx=10, pady=5)
entry5 = tk.Entry(root)
entry5.grid(row=4, column=1, padx=10, pady=5)
entry6 = tk.Entry(root)
entry6.grid(row=5, column=1, padx=10, pady=5)

# 创建按键
def show():
    #name = "\n"+entry1.get()+"~"
    ip = entry2.get()+"~"
    user = entry3.get()+"~"
    passwd = entry4.get()+"~"
    port = entry5.get()
    firename =  entry6.get()
    #a = name+ip+user+passwd+port
    #pwd = os.path.dirname(os.path.abspath(__file__))
    pwd = os.getcwd()
    pwffire = os.path.join(pwd,firename)
    with open(pwffire, 'a+',encoding="utf-8") as f:
        size = os.path.getsize(pwffire) 
        if size == 0:
            name = entry1.get()+"~"
            a = name+ip+user+passwd+port
        else:
            name = "\n"+entry1.get()+"~"
            a = name+ip+user+passwd+port
        f.write(a)


 



def run():
    #pwd = os.path.dirname(os.path.abspath(__file__))
    pwd = os.getcwd()
    pwffire = os.path.join(pwd,"Golin.exe -run linux")
    os.system(pwffire)
    tk.messagebox.showinfo('提示','采集完成！')

button1 = tk.Button(root, text='增加服务器信息', command=show).grid(row=6, column=0,
                                            sticky=tk.W, padx=30, pady=5)                                 
button2 = tk.Button(root, text='退出', command=root.quit).grid(row=6, column=1,
                                          sticky=tk.E, padx=30, pady=5)
button3 = tk.Button(root, text='运行', command=run).grid(row=6, column=2,
                                            sticky=tk.W, padx=30, pady=5) 

tk.mainloop()


