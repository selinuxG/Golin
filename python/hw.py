# -*- coding=utf-8 -*-
import paramiko
import time
import sys
import os


class ssh_con():
    def __init__(self, host, port, username, password):
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(hostname=host, port=port, username=username, password=password, allow_agent=False,
                    look_for_keys=False, timeout=10)
        self.command = ssh.invoke_shell()
        self.date = ""

    def cmd_date(self, cmd):
        for i in cmd:
            self.command.send(i)
            self.command.send("\n                      \n")
            time.sleep(2)
            self.output = self.command.recv(65535).decode('ascii')
            self.date += self.output
        return self.date


file_name = sys.argv[2] + "_" + sys.argv[3] + ".log"

file_date = ssh_con(host=sys.argv[3], port=sys.argv[6], username=sys.argv[4], password=sys.argv[5]).cmd_date(
    sys.argv[7].split(";"))

with open(os.path.join(sys.argv[1], file_name), "w+", encoding="utf-8") as f:
    for i in file_date:
        f.write(i.replace("\n", ""))
