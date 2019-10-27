package db

import (
  "time"
  "errors"
)

// Item : here you tell us what Item is
type Item interface {
  DataType() string
  ttl() int
}

// StringItem : here you tell us what StringItem is
type StringItem struct {
  Value string
  expiresAt time.Time
}

// Returns data type
func (si StringItem) DataType() string {
    return "string"
}

// Returns the number of seconds till the item expires, -1 if item doesn't expire
func (si StringItem) ttl() int {
  return -1
}

// setter
func (si *StringItem) Set(s string) {
  si.Value = s
}

// Getter
func (si StringItem) Get() string {
  return si.Value
}

// ListItem : here you tell us what ListItem is
type ListItem struct {
  Value []string
  expiresAt time.Time
}

// Returns data type
func (li ListItem) DataType() string {
    return "list"
}

// Returns the number of seconds till the item expires, -1 if item doesn't expire
func (li ListItem) ttl() int {
  return -1
}


/*
func (li ListItem) get() []string {
  return li.value
}
*/

// llen
func (li ListItem) llen() int {
  lenght := len(li.Value)
  return lenght
}



// Db : here you tell us what Db is
type Db map[string]Item

// implement set() and get() methods for Db type
func (di Db) Set(key string, i1 Item) {
  di[key] = i1
}


func (di Db) Get(key string) (Item,error) {
  if it, ok := di[key]; ok {
    return it,nil
  } 
  
  return nil,errors.New("Item does not exist")
  
}

