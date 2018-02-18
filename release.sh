NAME="freshreadme"
VERSION="0.5"

rm -R out/*

for TARGET in "linux" "darwin" "windows"
do
    DIRNAME="$NAME-$TARGET-$VERSION"
    BIN=$(test "$TARGET" = "windows" && echo "$NAME.exe" || echo "$NAME")
    env GOOS="$TARGET" GOARCH=amd64 go build -o "out/$DIRNAME/$BIN"
    cp README.md "out/$DIRNAME/"
    cp LICENSE.txt "out/$DIRNAME/"
    tar czvf "out/$DIRNAME.tar.gz" -C "out/" "$DIRNAME"
done
