If, at the end of a function (inside it),
```go
defer x.Something();
return y.Something();
```
then y.Something() runs first
then runs x.Something()
and then function returns with whatever y.Something() returned
