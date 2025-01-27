.class public Fib
.super java/lang/Object

.method public static main([Ljava/lang/String;)V
    .limit locals 1
    .limit stack 1
    getstatic java/lang/System/out Ljava/io/PrintStream;
    ldc 5
    invokestatic Fib/Fib(I)I
    invokevirtual java/io/PrintStream/println(I)V

    ldc 10
    invokestatic Fib/Fib(I)I
    invokevirtual java/io/PrintStream/println(I)V
    return
.end method

.method public static Fib(I)I
    .limit locals 1
    .limit stack 20
    ldc 0
    iload_0
    if_icmpgt ELSE_0
        iload_0
        ldc 2
        if_icmpgt ELSE_1
            ldc 1
            ireturn
            goto END_IF_ELSE_1
        ELSE_1:
            ldc 1
            iload_0
            isub
            invokestatic Fib/Fib(I)I

            ldc 2
            iload_0
            isub
            invokestatic Fib/Fib(I)I

            iadd
            ireturn
        END_IF_ELSE_1:
        goto END_IF_ELSE_0:
    ELSE_0:
    END_IF_ELSE_0:
    ldc 1
    ireturn
.end method