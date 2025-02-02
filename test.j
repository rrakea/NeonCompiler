.class public Test
.super java/lang/Object 

.field public static i I


.method static <clinit>()V 
.limit stack 1
ldc 2
putstatic Test/i I
return 
.end method




.method public static add(II)I
.limit stack 3
.limit locals 2
iload_0
iload_1
iadd
getstatic Test/i I
iadd
ireturn
ldc 0
ireturn
.end method


.method public static main([Ljava/lang/String;)V
.limit stack 7
.limit locals 2
ldc 0
istore_1
WHILE_BEGIN_0:
iload_1
ldc 10
ldc 0
invokestatic Test/add(II)I
if_icmpne IS_TRUE_0
ldc 0
goto BOOL_EX_END_0
IS_TRUE_0:
ldc 1
BOOL_EX_END_0:

ifeq WHILE_END_0
getstatic java/lang/System/out Ljava/io/PrintStream;
iload_1
invokevirtual java/io/PrintStream/println(I)V
iload_1
ldc 1
iadd
istore_1
goto WHILE_BEGIN_0
WHILE_END_0:
return
.end method

