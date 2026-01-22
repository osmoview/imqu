# ImQu

Provide format to write and read data sequence/queues from the files.

Format description:

```
[  8byte  ][   8byte   ][   data   ]...
[ wait ns ][ data size ][   data   ]...
```

At the end, of the file contains 8 bytes with count of frames in this file

## Usage

Write file:

```go
iq, err := imqu.NewImQuWriter(fn)
...

iq.AddData(b)
// or
iq.AddDataWithTime(b, time)
```

Read file:

```go
iq, err := imqu.OpenImQu(fn)
...

for {
    b, err := iq.NextFrame(1) // where 1 is speed coefficient 
    if err != nil{
        // check on the io.EOF
        break
    }

    ...
}
```
