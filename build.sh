mkdir -p build

go build -o build/tt cmd/*.go

if [ -d build/test ]
then
    cp build/tt build/test/tt
fi