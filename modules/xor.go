package modules

import "strings"

// initialize add to waiting list
// till load method gets invoked
func init() {
  loadIntenralModule("xor", XORHandler, "")
}

func XORHandler(command string, stdin string, args ...string) ([]byte, error){
  if stdin == "" {
    return []byte{}, nil
  }
  result := strings.SplitN(stdin, ":", 2)
  if len(result) < 2 {
    return []byte{}, nil        
  }      
  key, rest := result[0], result[1:]
  content := ""
  for _, item := range rest {
    content += item
  }

  var output []byte 
  for i := 0; i < len(content); i++ {
    output = append(output, content[i] ^ key[i % len(key)])
  }

  return output, nil
}
