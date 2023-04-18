# -*- coding:utf-8 -*-
import sys
if __name__ == '__main__':
    print(sys.argv)
    with open("test.txt","w",encoding="utf-8") as f:
        for i in sys.argv:
            f.write(i+"\n")