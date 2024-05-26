package main

import (
	"fmt"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib"
)

func main() {
	kv := NewKVStore()

	fmt.Println(kv.Set("foo", "bar"))    
	fmt.Println(kv.Get("foo"))           
	fmt.Println(kv.Strlen("foo"))        
	fmt.Println(kv.Append("foo", "bazing")) 
	fmt.Println(kv.Get("foo"))           
	fmt.Println(kv.Delete("foo"))        
	fmt.Println(kv.Get("foo"))
}
