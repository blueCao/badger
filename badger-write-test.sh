#/bin/bash#
#
#
#
#Date
#	2018-1-3
#
#Author
#	Ove
#
#测试badger写入性能的脚本，需要输入3个参数:写入多少条数据   最大的线程数量    写入的字符串长度单位为byte      数据库文件写入的目录
#

#测试结果目录、数据库文件目录
result_dir="badger_test_result"
badger_dir=$4
#检查是否有这个文件夹,没有的话创建
if [ ! -d "$result_dir" ]; then
  mkdir $result_dir
fi
if [ ! -d "$badger_dir" ]; then
  mkdir $badger_dir
fi

#检查参数
if [ "$#" -lt  4 ]; then
  echo "输入的参数个数不够，输入参数个数有误需要输入4个，  写入多少条数据   最大的线程数量    写入的字符串长度单位为byte     数据库文件写入的目录    "
fi
echo "执行写性能测试, 写入 $1 条数据，线程$2个，写入的字符串长度 $3 byte   写入目录$4" > $result_dir/badger-write-test-$1-$2-$3.result

#显示数据库原始目录大小
echo "数据库写入前目录大小为" >> $result_dir/badger-write-test-$1-$2-$3.result
du -sh $badger_dir  >> $result_dir/badger-write-test-$1-$2-$3.result 

#开始时间写入结果中i
echo "开始时间：" >> $result_dir/badger-write-test-$1-$2-$3.result
date >> $result_dir/badger-write-test-$1-$2-$3.result

#执行测试脚本 "$*"表示使用所有的参数（注意使用"$*"必须加上“”双引号）
#记录错误日志
./badger-write-test $*  1>/dev/null  2>$result_dir/badger-write-test-$1-$2-$3.error

#结束时间写入结果中
echo "结束时间：" >> $result_dir/badger-write-test-$1-$2-$3.result
date >> $result_dir/badger-write-test-$1-$2-$3.result

#错误日志信息为空则删除日志文件
error_log_size=$(stat -c%s "$result_dir/badger-write-test-$1-$2-$3.error")
if [ $error_log_size -eq 0 ]; then
  rm $result_dir/badger-write-test-$1-$2-$3.error
fi

#显示数据库目录大小
echo "数据库写入后目录大小为"  >> $result_dir/badger-write-test-$1-$2-$3.result
du -sh $badger_dir >> $result_dir/badger-write-test-$1-$2-$3.result 
