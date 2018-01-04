#/bin/bash#
#
#
#
#Date
#	2018-1-4
#
#Author
#	Ove
#
#测试badger查询性能的脚本，需要输入2个参数：随机查询多少条数据   同时启动的线程数量
#

#测试结果目录、数据库文件目录
result_dir="badger_test_result"
badger_dir=/home/caojunhui/workspace/go/badger
#检查是否有这个文件夹,没有的话创建
if [ ! -d "$result_dir" ]; then
  mkdir $result_dir
fi

#检查参数
if [ "$#" -lt  2 ]; then
  echo "输入的参数个数不够，需要输入2个参数：随机查询多少条数据   同时启动的线程数量"
fi
echo "执行badger查询性能测试, 随机查询 $1 条数据，启动查询线程$2个" > $result_dir/badger-query-test-$1-$2.result

#显示当前数据库目录大小
badger_dir=/home/caojunhui/workspace/go/badger
echo "数据库查询前目录大小为" >> $result_dir/badger-query-test-$1-$2.result
du -sh $badger_dir  >> $result_dir/badger-query-test-$1-$2.result 

#开始时间写入结果中
echo "开始时间：" >> $result_dir/badger-query-test-$1-$2.result
date >> $result_dir/badger-query-test-$1-$2.result

#执行测试脚本 "$*"表示使用所有的参数（注意使用"$*"必须加上“”双引号）
#记录错误日志
./badger-query-test $*  1>/dev/null  2>$result_dir/badger-query-test-$1-$2.error

#结束时间写入结果中
echo "结束时间：" >> $result_dir/badger-query-test-$1-$2.result
date >> $result_dir/badger-query-test-$1-$2.result

#错误日志信息为空则删除日志文件
error_log_size=$(stat -c%s "$result_dir/badger-query-test-$1-$2.error")
if [ $error_log_size -eq 0 ]; then
  rm $result_dir/badger-query-test-$1-$2.error
fi

#显示数据库目录大小
echo "数据库查询后目录大小为"  >> $result_dir/badger-query-test-$1-$2.result
du -sh $badger_dir >> $result_dir/badger-query-test-$1-$2.result 
