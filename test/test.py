#!/usr/bin/python

import sys
import time
import logging
from logging.handlers import RotatingFileHandler

console1 = logging.StreamHandler(sys.stderr)
console1.setLevel(logging.ERROR)
console2 = logging.StreamHandler(sys.stdout)
console2.setLevel(logging.DEBUG)
logging.basicConfig(level=logging.DEBUG,
                format='%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s',
                datefmt='%a, %d %b %Y %H:%M:%S',
                handlers=[console1, console2])

'''
logging.basicConfig(level=logging.DEBUG,
                format='%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s',
                datefmt='%a, %d %b %Y %H:%M:%S',
                filename='test_log.log',
                filemode='w')
'''

#################################################################################################
#定义一个StreamHandler，将INFO级别或更高的日志信息打印到标准错误，并将其添加到当前的日志处理对象#
'''
console = logging.StreamHandler(sys.stderr)
console.setLevel(logging.INFO)
formatter = logging.Formatter('%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s', '%a, %d %b %Y %H:%M:%S')
console.setFormatter(formatter)
logging.getLogger('').addHandler(console)
'''
#################################################################################################


#################################################################################################
#定义一个RotatingFileHandler，最多备份50个日志文件，每个日志文件最大10M
'''
Rthandler = RotatingFileHandler('test_log.log', maxBytes=10*1024*1024,backupCount=50)
Rthandler.setLevel(logging.DEBUG)
formatter = logging.Formatter('%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s', '%a, %d %b %Y %H:%M:%S')
Rthandler.setFormatter(formatter)
logging.getLogger('').addHandler(Rthandler)
'''
################################################################################################

logging.info("Start : %s", time.ctime())
logging.debug("Debug....")

while  True:
    time.sleep( 60 )
    logging.error("NOW : %s", time.ctime())
    logging.debug("Debug....")
