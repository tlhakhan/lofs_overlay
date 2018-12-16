package main

import(
  "bytes"
  "io/ioutil"
  "bufio"
  "os/exec"
  "strings"
  "path/filepath"
  "flag"
  "fmt"
  "os"
)


func main() {
  overlayRootPath := flag.String("overlayRootPath", "/usbkey/crud", "The folder path that will be considered as / when overlaying")
  flag.Parse()

  // expecting last argument to be "start" or "stop"
  switch(flag.Arg(0)) {
    case "start":
      fmt.Println("start the overlay process")
      start(*overlayRootPath)
      fmt.Println("done linking the overlay")
    case "stop":
      fmt.Println("stop the overlay process")
      stop(*overlayRootPath)
      fmt.Println("done unlinking the overlay")
    default:
      fmt.Printf("Usage: %s start|stop\n", string(os.Args[0]))
      os.Exit(1)
    }
}

/*
  Takes in an argument that is considered as the path that will be overlayed on / and then unmounts the files accordingly

  - scan the /etc/mnttab for the files that correspond to mount that the beginning of the overlayRootPath
  - place the files in an array
  - then process each file with umount, quit out on any unmount error
*/
func stop(overlayRootPath string) {
  // prior to mounting check if the targetPath is already mounted
  m, err := ioutil.ReadFile("/etc/mnttab")
  if err != nil {
    fmt.Printf("error reading /etc/mnttab, exiting..\n")
    fmt.Println(err)
    os.Exit(1)
  }
  scanner := bufio.NewScanner(bytes.NewReader(m))
  targetList := []string{}
  for scanner.Scan() {
    line := scanner.Text()
    fields := strings.Fields(line)
    if strings.HasPrefix(fields[0], overlayRootPath){
      targetList = append(targetList, fields[1])
    }
  }
  fmt.Printf("unmounting the found target files: %q\n", targetList)
  for _, targetPath := range targetList {
    _, err:= exec.Command("umount", targetPath).CombinedOutput()
    if err != nil {
      fmt.Printf("error unmounting target file %s, exiting..", targetPath)
      os.Exit(1)
    }
    fmt.Printf("target file %s is now unmounted\n", targetPath)
  }
}

/*
  Takes in an argument that is considered as the path that will be overlayed on /
  
  - make sure that the overlayRootPath exists and is a directory
  - traverse into the directory
    - for every directory
      - if the directory exists at target location then do nothing return nil
      - if the directory doesn't exist at target location then create it
      - if the directory exists, but is a file, then quit out for manual intervention
    - for every file
      - if the file exists at target location, then continue
      - if the file doesn't exist then create an empty file, then continue
      - if the file exists, but is a directory, then quit out for manual intervention
      - before mounting over, check if file is already overlayed in /etc/mnttab
        - if file is not overlayed, then overlay the file over the target location
        - if file is already overlayed, then do nothing
*/
func start(overlayRootPath string) {
  // ensure that the rootPath exists
  overlayRootPathInfo, err := os.Stat(overlayRootPath)
  if os.IsNotExist(err) {
    fmt.Printf("%s does not exist, exiting..\n", overlayRootPath)
    os.Exit(1)
  }
  if !overlayRootPathInfo.IsDir() {
    fmt.Printf("%s is not a directory, exiting..\n", overlayRootPath)
    os.Exit(1)
  }
  // walk the folder structure and create accordingly on each walk
  fmt.Printf("walking dir %s\n", overlayRootPath)
  filepath.Walk(overlayRootPath, func(path string, info os.FileInfo, err error) error {
    if info.IsDir(){
      if path == overlayRootPath {
        // ignore processing the root walkpath directory
        return nil
      }
      fmt.Printf("found directory %s\n", path)
      // convert overlay source path to target path
      targetPath := strings.TrimPrefix(path, overlayRootPath)
      targetFileInfo, err := os.Stat(targetPath)
      // the dir can exist
      // the dir needs to be created
      // something there but not dir
      if !os.IsNotExist(err) {
        fmt.Printf("target path %s exists\n", targetPath)
        if targetFileInfo.IsDir() {
          // no work needs to be done
          fmt.Printf("target path %s is a directory\n", targetPath)
          fmt.Printf("mkdir is not needed on %s\n", targetPath)
          return nil
        } else {
          // targetPath is a file instead of dir
          fmt.Printf("%s is a dir, but target path %s is a file, exiting..\n", path, targetPath)
          os.Exit(1)
        }
      } else {
        // dir needs to be made
        if err := os.Mkdir(targetPath, 755); err != nil {
          fmt.Printf("error when creating target path %s, exiting..\n", targetPath)
          fmt.Println(err)
          os.Exit(1)
        }
      }
      return nil
    } else {
      fmt.Printf("found file %s\n", path)
      targetPath :=strings.TrimPrefix(path, overlayRootPath)
      targetFileInfo, err := os.Stat(targetPath)

      if !os.IsNotExist(err) {
        // file exists
        fmt.Printf("target path %s exists\n", targetPath)
        if targetFileInfo.IsDir() {
          fmt.Printf("%s is a file, but target path %s is a directory, exiting..\n", path, targetPath)
          os.Exit(1)
        } else {
          fmt.Printf("target path %s is a file\n", targetPath)
          fmt.Printf("create file is not needed on %s\n", targetPath)
        }
      } else {
        // file needs to be touched
        _, err := os.Create(targetPath)
        if err != nil {
          fmt.Printf("error creating target file %s", targetPath)
          fmt.Println(err)
          os.Exit(1)
        }
      }
      // prior to mounting check if the targetPath is already mounted
      m, err := ioutil.ReadFile("/etc/mnttab")
      if err != nil {
        fmt.Printf("error reading /etc/mnttab, exiting..\n")
        fmt.Println(err)
        os.Exit(1)
      }
      scanner := bufio.NewScanner(bytes.NewReader(m))
      found := false
      for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Fields(line)
        if fields[0] == path && fields[1] == targetPath {
          found = true
          break
        }
      }

      if !found {
        //  mount -F lofs "$file" "$dest_file"
        // mount happens here, all edge cases have errored out
        _, err:= exec.Command("mount", "-F", "lofs", path, targetPath).CombinedOutput()
        if err != nil {
          fmt.Printf("error mounting lofs for %s to target file %s", path, targetPath)
          os.Exit(1)
        }
        fmt.Printf("%s is now mounted to target file %s\n", path, targetPath)
      } else {
        fmt.Printf("%s was already mounted to target file %s\n", path, targetPath)
      }
      return nil
    }
  })
}
