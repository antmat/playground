# playground
filesorting tool written in go, using external sort

## Building
  
```
go get github.com/antmat/playground/generator
go get github.com/antmat/playground/sorter
```
  
## Usage
```
$GOPATH/bin/generator {out_file} {line_count} {line_length}
$GOPATH/bin/sorter {infile} {outfile} {tmp_dir} {tmp_file_size}
```

## Known bugs and limitations
  * No tests
  * Files should contain ASCII text '\n' divided
  * Huge lines without '\n' delimiters can result in OOM
  * tmp_file_size is not determined atomatically and is not correctly limited:
    * huge value can result in OOM
    * small values will result in degraded performance or event resource (fd) exhausting 
  * Interrupted execution can leave temporary files on fs
  * All output files are recreated (e.g. symbolyc links are replaced with normal files)
  * No command line flag parser