.class public AddTest
.super java/lang/Object 


.method public static add(II)I
.limit stack 2
.limit locals 2
iload_0
iload_1
iadd
ireturn
.end method

.method public static main([Ljava/lang/String;)V
.limit stack 5
.limit locals 3
ldc 5
istore_1
ldc 10
istore_2
iload_1

iload_2

invokestatic AddTest/add(II)I
return 
.end method

