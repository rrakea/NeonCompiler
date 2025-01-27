# NEON

Compiler from a subset of C# to JVM Jasmin Bytecode 

For the compiler course at HHU Düsseldorf

By Konrad Burgi

## How to run and compile

Build:
go build neon.go

Run:
./main -compile [filepath]

Build Jasmin:
java -jar [path to jasmin.jar] [filepath]

Flags:
-liveness for variable liveness analysis
-constant for constant propogation analysis

## Info

Uses golang 1.23.2 & Jasmin 2.4
