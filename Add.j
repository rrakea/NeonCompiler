.class public Add
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
    .limit stack 3
    .limit locals 3
    ldc 5
    ldc 10
    invokestatic Add/add(II)I
    istore_1
    getstatic java/lang/System/out Ljava/io/PrintStream;
    iload_1
    invokevirtual java/io/PrintStream/println(I)V
    return
.end method

