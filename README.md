
# NEON

Compiler from a subset of C# to JVM Bytecode (Jasmin)

For the course compiler construction at HHU Düsseldorf

By Konrad Burgi

## How to run and compile

Build:
go build neon.go

Run:
./main -compile [filepath]

-liveness for variable liveness analysis
-constant for constant propogation analysis

## Info

Uses go 1.23.2 & 2.4
