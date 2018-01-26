#/bin/bash
#
#依次执行每一个查询，综合测试badger读性能

#查询数量
#query_records=1

#并发查询的线程数
#query_threads=1

#查询的范围【1，1200000000】
#query_scope=1200000000
query_scope=1000000000

#记录查询测试的进度日志
schedule_log=badger-query-test-all.log

#触发查询测试的执行脚本文件
query_sh=badger-query-test.sh

for query_threads in {256,128,64,512}
do
	for query_records in {1,100,1000,10000,100000,1000000,10000000}
	do
		#记录脚本启动的时刻
		echo "开始：	$query_records		$query_threads		$query_scope"	>>	$schedule_log
		date >> $schedule_log
		echo "sh $query_sh $query_records $query_threads $query_scope"
		sh $query_sh $query_records $query_threads $query_scope
		#记录脚本结束时刻
		echo "结束：    $query_records          $query_threads          $query_scope"   >>      $schedule_log
		date >> $schedule_log
		#换行
		echo "" >> $schedule_log	
	done
done

echo "所有查询测试完成！请查看日志确认无误	"$schedule_log
