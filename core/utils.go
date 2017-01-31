package core

import (
  "io/ioutil"
  "os"
  "bufio"
  "errors"
)

func FileExist(path string) (bool, error) {
  _, err := os.Stat(path)
  if err == nil {
    return true, nil
  } else if os.IsNotExist(err) {
    return false, err
  }
  return true, nil  
}

func GetPipedStdin() (string, error){
  fi, _ := os.Stdin.Stat()

  if (fi.Mode() & os.ModeCharDevice) == 0 {
    bytes, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
      return "", err
    }
    str := string(bytes)
    return str, nil
  }
  return "", errors.New("no piped stdin")
}

func GetTermStdin() (string, error) {
  fi, _ := os.Stdin.Stat()
  if (fi.Mode() & os.ModeCharDevice) != 0 {    
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
      return "", err
    }
    return input, nil
  }
  return "", errors.New("no stdin from term")
}

func SelectStdin() (string, error) {
  pipeout, err := GetPipedStdin()
  if err == nil {
    return pipeout, nil
  }
  termout, err := GetTermStdin()
  if err == nil {
    return termout, nil
  }
  return "", errors.New("no data in stdin")
}
