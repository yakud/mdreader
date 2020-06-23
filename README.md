# MD Reader

## first install snappy 
https://github.com/andrix/python-snappy

### mac OSX:
> brew install snappy
>
>$ CPPFLAGS="-I/usr/local/include -L/usr/local/lib" pip3 install python-snappy

## sencond make uncompressed data
> python3 -m snappy -d BTCUSD.dat BTCUSD.dat.uncompressed

## third make csv:
> go run ./cmd/mdreader -in=./BTCUSD.dat.uncompressed -out=./BTCUSD.csv
