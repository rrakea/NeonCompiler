.class public If
.super java/lang/Object 



.method public static main([Ljava/lang/String;)V
.limit stack 5
.limit locals 1
ldc 1
ldc 1
ifeq BOOL_EX_FALSE_0
ldc 1
goto BOOL_EX_END_0
BOOL_EX_FALSE_0:
ldc 0
BOOL_EX_END_0:

ldc 0
if_icmpeq ELSE_LABEL_0
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc "Passed"
invokevirtual java/io/PrintStream/println(Ljava/lang/String;)V
return
goto END_IF_ELSE_0
ELSE_LABEL_0
END_IF_ELSE_0:
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc "Not passed"
invokevirtual java/io/PrintStream/println(Ljava/lang/String;)V
return 
.end method

