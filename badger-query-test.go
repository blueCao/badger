/*
 * @Author Ove
 * @Date 2018-1-4
 *
 * 数据库目录位于： /home/caojunhui/workspace/go/badger
 * 测试badger key-value ssd 数据库的查询性能
 * 带3个参数（空格分隔）：随机查询多少条数据   启动同时查询的线程数量      查询的区间
 *
 */
package main

import (
       "math/rand"
        "strconv"
        "runtime"
        "os"
        "log"
        "fmt"
        "github.com/dgraph-io/badger"
)

//query a specific value from database by key
func query_job(db *badger.DB, query_chan <-chan []byte, exit_signal chan<- bool) {
  for{
      //get an query key
      key :=<- query_chan
      //the last finished signal "0"
      if string(key[:]) == "0" {
          //send finshed signal
          exit_signal <- true
      }
      //query
      db.View(func(txn *badger.Txn) error {
        item, err := txn.Get(key)
        if err != nil {
          fmt.Printf("查询key：%s \n",string(key[:]))
          //send stop signal to main thread
          exit_signal <- false
          return nil
        }
        val, err := item.Value()
        if err != nil {
          fmt.Printf("查看key：%s value：%s 出错\n",string(key[:]), val)
          //send stop signal to main thread
          exit_signal <- false
          return nil
        }
        fmt.Printf("key: %s value：%s \n", string(key[:]), val)
        return nil
      })
    }
}


//randam query input routine
func input_job(size int, query_chan chan<- []byte, scope int){
  var random_value int = -1
  for i := 1; i <= size; i++ {
      //get random num between 1 and size
      random_value = rand.Intn(scope) + 1
      //convert the int into []byte
      query_chan <- []byte(strconv.Itoa(random_value))
      fmt.Printf("输入随机数： %d\n", random_value)
  }
  //the last finished signal "0"
  query_chan <- []byte(strconv.Itoa(0))
  fmt.Println("随机数输入完成")
}

func main() {
  // Open the Badger database located in the /tmp/badger directory.
  // It will be created if it doesn't exist.
  opts := badger.DefaultOptions
  opts.Dir = "/home/caojunhui/workspace/go/badger"
  opts.ValueDir = "/home/caojunhui/workspace/go/badger"
  db, err := badger.Open(opts)
  if err != nil {
          fmt.Println("创建数据库文件错误"," ",opts.Dir," ",opts.ValueDir)
          log.Fatal(err)
  }
  defer db.Close()

  //get input arguments as query size records num and GOMAXPROCS
  args := os.Args[1:]
  if len(args) < 3 {
      fmt.Println("输入参数个数有误 3个参数（空格分隔）：随机查询多少条数据   启动同时查询的线程数量  查询的区间")
      return
  }
  var size, gomaxprocs,scope int = 0, 0, 0
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
    scope, err = strconv.Atoi(args[2])
    if err != nil || scope <= 0{
        fmt.Println("输入参数 ",args[2]," 有误，请重新输入正整数")
        return
    }
    fmt.Printf("随机查询 %d 多少条数据,  同时启动 %d 个线程数量   查询key值区间【1，%d】之间的value \n",size, gomaxprocs,scope)


  //set the max procs
  runtime.GOMAXPROCS(gomaxprocs)

  //start a random query input thread
  keys_chan := make(chan []byte, gomaxprocs);
  exit_signal := make(chan bool)
  go input_job(size, keys_chan,scope)

  //start multi write threads
  for i := 0; i < gomaxprocs; i++  {
      go query_job(db, keys_chan, exit_signal)
  }

  //waiting exit signal
  signal :=<- exit_signal
  if signal==false {
      fmt.Println("查询数据异常，程序退出！")
  } else {
        fmt.Println("全部查询完成!")
  }

}
