package main

import (
  "time"
  //"fmt"
  "errors"
)

// Item : here you tell us what Item is
type Item interface {
  dataType() string
  ttl() int
}

// StringItem : here you tell us what StringItem is
type StringItem struct {
  value string
  expiresAt time.Time
}

// Returns data type
func (si StringItem) dataType() string {
    return "string"
}

// Returns the number of seconds till the item expires, -1 if item doesn't expire
func (si StringItem) ttl() int {
  return -1
}

// setter
func (si *StringItem) set(s string) {
  si.value = s
}

// getter
func (si StringItem) get() string {
  return si.value
}

// ListItem : here you tell us what ListItem is
type ListItem struct {
  value []string
  expiresAt time.Time
}

// Returns data type
func (li ListItem) dataType() string {
    return "list"
}

// Returns the number of seconds till the item expires, -1 if item doesn't expire
func (li ListItem) ttl() int {
  return -1
}

// rpush
func (li *ListItem) rpush(s string) {
  li.value = append(li.value, s)
}

// rpop
func (li ListItem) rpop() string {
  last := li.value[len(li.value)-1]
  return last
}

func (li ListItem) get() []string {
  return li.value
}

// llen
func (li ListItem) llen() int {
  lenght := len(li.value)
  return lenght
}



// Db : here you tell us what Db is
type Db map[string]Item

// implement set() and get() methods for Db type

func (di Db) set(key string, i1 Item) {
  di[key] = i1
  //fmt.Println(di[key])
}


func (di Db) get(key string) (Item,error) {
  if it, ok := di[key]; ok {
    return it,nil
  } 
  
  return nil,errors.New("Item does not exist")
  
}

