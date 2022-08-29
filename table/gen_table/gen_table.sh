#!/bin/bash

CURDIR=$(cd $(dirname $0); pwd)
TEMPLATE=$CURDIR/table_template.txt
TABLE_LIST=$CURDIR/table_list.txt
OUT=$CURDIR/..

CMD="gen_table.exe"
if [ $(uname -s) = 'Linux' ]; then
	CMD="gen_table_linux"
elif [ $(uname -s) = 'Darwin' ]; then
    CMD="gen_table_mac"
fi

$CURDIR/$CMD -template=$TEMPLATE -tablelist=$TABLE_LIST -out=$OUT
read -p "按任意键继续" -n 1 -r

