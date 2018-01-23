/*
 * @Author Ove
 * @Date 2018-1-4
 *
 * 数据库目录位于： /home/caojunhui/workspace/go/badger
 * 测试badger key-value ssd 数据库的查询性能
 * 带2个参数（空格分隔）：随机查询多少条数据   启动同时查询的线程数量
 *
 */
package main

import (
    "log"
    "github.com/dgraph-io/badger"
    "fmt"
)

func main() {
  // Open the Badger database located in the /tmp/badger directory.
  // It will be created if it doesn't exist.
  opts := badger.DefaultOptions
  opts.Dir = "/Users/Ove/workspace/go/badger/badgerDB"
  opts.ValueDir = "/Users/Ove/workspace/go/badger/badgerDB"
  db, err := badger.Open(opts)
  if err != nil {
          fmt.Println("创建数据库文件错误"," ",opts.Dir," ",opts.ValueDir)
          log.Fatal(err)
  }
  defer db.Close()

    //update(write)
    err = db.Update(func(txn *badger.Txn) error {
      err := txn.Set([]byte("name"), []byte("my name is Ove"))
      return err
    })
    var query_result []byte
    //read
    err = db.View(func(txn *badger.Txn) error {
      item, err := txn.Get([]byte("name"))
      if err != nil {
        return err
      }
      val, err := item.Value()
      if err != nil {
        return err
      }
      query_result = val
      return nil
    })
    if err == nil {
        fmt.Printf("%s \n",query_result)
    }
    //update(delete) the same as update(write)
    err = db.Update(func(txn *badger.Txn) error {
      err := txn.Delete([]byte("name"))
      return err
    })
    //read empty
    err = db.View(func(txn *badger.Txn) error {
      item, err := txn.Get([]byte("name"))
      if err != nil {
        return err
      }
      val, err := item.Value()
      if err != nil {
        return err
      }
      query_result = val
      return nil
    })
    if err != nil && err == badger.ErrKeyNotFound {
      fmt.Println("key not found")
    }

}
