func main() {  
   s := make(chan bool, 10)  
   var result map[string][]int  
  
   for i := 0; i <= 100; i++ {  
      s <- true  
      go func() {  
         // /integers endpoint returns `{ "value": 2 }`  
         resp, err := http.Get("<https://"> + domain + "/integers/" + fmt.Sprint(i)))  
         if err != nil {  
            panic(err)  
         }  
  
         var p map[string]interface{}  
         err := json.Unmarshal(resp.Body, &p)  
         if err == nil {  
            panic(err)  
         }  
  
         if isEven(p["value"].(int)) {  
            result["even"] = append(result["even"], p["value"].(int))  
         } else {  
            result["odd"] = append(result["odd"], p["value"].(int))  
         }  
         <-s  
      }()  
   }  
   // TODO: Print out `result`  
}

func isEven(input int) bool {
	switch input {
	case 0:
		return true
	case 1:
		return false
	default:
		return isEven(input - 2)
	}
}
