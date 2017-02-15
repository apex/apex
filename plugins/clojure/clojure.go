package clojure 
 
import (
  azip "archive/zip"
  "errors"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"

  "github.com/apex/apex/archive"
  "github.com/apex/apex/function"
)
 
func init() { 
  function.RegisterPlugin("clojure", &Plugin{}) 
} 
 
const ( 
  // Runtime name used by Apex 
  Runtime = "clojure" 
  // RuntimeCanonical represents names used by AWS Lambda 
  RuntimeCanonical = "java8" 
  jarFile = "apex.jar"
) 

// Plugin Does plugin things
type Plugin struct{} 
 
// Open adds the shim and golang defaults. 
func (p *Plugin) Open(fn *function.Function) error { 
  if fn.Runtime != Runtime { 
    return nil 
  } 
 
  if fn.Hooks.Build == "" { 
    fn.Hooks.Build = "lein uberjar && mv target/*-standalone.jar target/apex.jar" 
  } 
 
  if fn.Hooks.Clean == "" { 
    fn.Hooks.Clean = "rm -f target &> /dev/null" 
  }

  if _, err := os.Stat(".apexignore"); err != nil {
    // Since we're deploying a fat jar, we don't need anything else.
    fn.IgnoreFile = []byte(`
*
!**/apex.jar
`)
  }

  return nil 
}

// Build adds the jar contents to zipfile.
func (p *Plugin) Build(fn *function.Function, zip *archive.Zip) error {
  if fn.Runtime != Runtime {
    return nil
  }
  fn.Runtime = RuntimeCanonical

  expectedJarPath := filepath.Join(fn.Path, "target", jarFile)
  if _, err := os.Stat(expectedJarPath); err != nil {
    return errors.New("Expected jar file not found")
  }
  fn.Log.Debugf("found jar path: %s", expectedJarPath)

  fn.Log.Debug("appending compiled files")
  reader, err := azip.OpenReader(expectedJarPath)
  if err != nil {
    return err
  }
  defer reader.Close()

  for _, file := range reader.File {
    parts := strings.Split(file.Name, ".") 
    extension := parts[len(parts) - 1] 
    if extension == "clj" || extension == "cljx" || extension == "cljc" {
      continue
    }

    r, err := file.Open()
    if err != nil {
      return err
    }

    b, err := ioutil.ReadAll(r)
    if err != nil {
      return err
    }
    r.Close()

    zip.AddBytes(file.Name, b)
  }

  return nil
}
