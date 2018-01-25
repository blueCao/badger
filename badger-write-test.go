/*
 * @Author Ove
 * @Date 2018-1-3
 *
 * 数据库目录位于： /home/caojunhui/workspace/go/badger
 * 测试badger key-value ssd 数据库的写入性能
 * 带3个参数（空格分隔）：写入多少条数据   最大的线程数量    写入的字符串长度单位为byte
 *
 * Badger Write-Benchmark
 *
 *
 * Three Input arguments：
 *          size            ：   how many records were writed into
 *          gomaxprocs      :   the gomaxprocs. Setting the max writing thread and setting runtime.GOMAXPROCS(gomaxprocs)
 *          value_length(bytes)    :   each value with specific length
 *          Badger DB Directory     : for example "/home/caojunhui/workspace/go/badger"
 */
package main

import (
        "bytes"
        "strconv"
        "runtime"
        "os"
        "log"
        "fmt"
        "github.com/dgraph-io/badger"
)

//write a k-v pair into database
func update(db *badger.DB, key []byte, value []byte) error {
  err := db.Update(func(txn *badger.Txn) error {
    err := txn.Set(key, value)
    return err 
  })
  return err
}

//input routine
func input_job(size int, data_chan chan<- []byte){
  for i := 1; i <= size; i++ {
      //convert the int into []byte
      data_chan <- []byte(strconv.Itoa(i))
      fmt.Printf("输入数值： %d\n", i)
  }
  fmt.Println("数据输入完成")
}

//write routine from channel
func write_job(db *badger.DB, input <-chan []byte, value []byte, exit_signal chan<- bool, last_key string) {
  for{
    //read an input key from chan
    key :=<- input 
    //write
    erro := update(db, key, value)
    if erro != nil {
      fmt.Printf("写入数据%s,%s出错\n ", string(key[:]), string(value[:]))
      fmt.Println(erro)
      //send stop signal to main thread
      exit_signal <- false
    } else {
	fmt.Printf("成功写入数据%s,%s\n ", string(key[:]), string(value[:]))
    }
    if string(key[:]) == last_key {
      //send finshed signal
      exit_signal <- true
    }
  }
}

func main() {
    //get input arguments as writed records num and GOMAXPROCS
    args := os.Args[1:]
    if len(args) < 4 {
        fmt.Println("输入参数个数有误需要输入4个，文件目录      写入多少条数据   最大的线程数量    写入的字符串长度单位为byte")
        return
    }

    // Open the Badger database located in the /tmp/badger directory.
  // It will be created if it doesn't exist.
  opts := badger.DefaultOptions
  opts.Dir = args[3]
  opts.ValueDir = args[3]
  db, err := badger.Open(opts)
  if err != nil {
          fmt.Println("创建数据库文件错误"," ",opts.Dir," ",opts.ValueDir)
          log.Fatal(err)
  }
  defer db.Close()

  var size, gomaxprocs, data_len int = 0, 0, 0
  size, err = strconv.Atoi(args[0])
  if err != nil || size <= 0{
      fmt.Println("输入参数 ",args[0]," 有误，请重新输入正整数")
      return
  }
  gomaxprocs, err = strconv.Atoi(args[1])
  if err != nil || gomaxprocs <= 0{
      fmt.Println("输入参数 ",args[1]," 有误，请重新输入正整数")
      return
  }
  data_len, err = strconv.Atoi(args[2])
  if err != nil || data_len <= 0{
      fmt.Println("输入参数 ",args[2]," 有误，请重新输入正整数")
      return
  }
  fmt.Printf("写入 %d 多少条数据,  最大 %d 个线程数量 ,  写入的字符串长度单位为 %d byte\n",size, gomaxprocs, data_len)

  //construct value all byte replace with '1'
  value := bytes.Repeat([]byte("1"), data_len)

  //set the max procs
  runtime.GOMAXPROCS(gomaxprocs)

  //start an input thread
  keys_chan := make(chan []byte, gomaxprocs);
  exit_signal := make(chan bool)
  go input_job(size, keys_chan)

  //start multi write threads
  for i := 0; i < gomaxprocs; i++  {
      go write_job(db, keys_chan, value, exit_signal,strconv.Itoa(size))
  }

  //waiting exit signal
  signal :=<- exit_signal
  if signal==false {
      fmt.Println("写入数据异常，程序退出！")
  } else {
        fmt.Println("全部写入完成!")
  }

}
