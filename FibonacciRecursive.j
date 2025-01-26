.class public FibonacciRecursive
.super java/lang/Object 



.method public static Fib(I)I
.limit stack 4
.limit locals 1
iload_0
ldc 0
if_icmpeq IS_TRUE_0
ldc 0
goto BOOL_EX_END_0
IS_TRUE_0:
ldc 1
BOOL_EX_END_0:

ifeq ELSE_LABEL_0
ldc 0
ireturn
goto END_IF_ELSE_0
ELSE_LABEL_0:
END_IF_ELSE_0:
iload_0
ldc 1
if_icmpeq IS_TRUE_1
ldc 0
goto BOOL_EX_END_1
IS_TRUE_1:
ldc 1
BOOL_EX_END_1:

ifeq ELSE_LABEL_1
ldc 1
ireturn
goto END_IF_ELSE_1
ELSE_LABEL_1:
END_IF_ELSE_1:
iload_0
ldc 1
isub
invokestatic FibonacciRecursive/Fib(I)I
iload_0
ldc 2
isub
invokestatic FibonacciRecursive/Fib(I)I
iadd
ireturn
ldc 0
ireturn
.end method


.method public static main([Ljava/lang/String;)V
.limit stack 2
.limit locals 1
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc 5
invokestatic FibonacciRecursive/Fib(I)I
invokevirtual java/io/PrintStream/println(I)V
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc 10
invokestatic FibonacciRecursive/Fib(I)I
invokevirtual java/io/PrintStream/println(I)V
return
.end method

