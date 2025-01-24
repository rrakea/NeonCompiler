.class public FibonacciRecursive
.super java/lang/Object 



.method public static Fib(I)I
.limit stack 7
.limit locals 1
iload_0
ldc 0
if_icmpgt BOOL_EX_FALSE_0
ldc 1
goto BOOL_EX_END_0
BOOL_EX_FALSE_0:
ldc 0
BOOL_EX_END_0:

ldc 0
if_icmpeq ELSE_LABEL_1
iload_0
ldc 2
ificmpge BOOL_EX_FALSE_1
ldc 1
goto BOOL_EX_END_1
BOOL_EX_FALSE_1:
ldc 0
BOOL_EX_END_1:

ldc 0
if_icmpeq ELSE_LABEL_0
ldc 1
ireturn
goto END_IF_ELSE_0
ELSE_LABEL_0END_IF_ELSE_0:
invokestatic FibonacciRecursive/Fib()I
invokestatic FibonacciRecursive/Fib()I
dadd
return
goto END_IF_ELSE_1
ELSE_LABEL_1END_IF_ELSE_1:
ldc 1
ireturn
.end method


.method public static main([Ljava/lang/String;)V
.limit stack 1
.limit locals 1
getstatic java/lang/System/out Ljava/io/PrintStream;
invokestatic FibonacciRecursive/Fib()I
invokevirtual java/io/PrintStream/println(I)V
getstatic java/lang/System/out Ljava/io/PrintStream;
invokestatic FibonacciRecursive/Fib()I
invokevirtual java/io/PrintStream/println(I)V
return 
.end method

