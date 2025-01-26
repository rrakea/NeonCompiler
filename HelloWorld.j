.class public HelloWorld
.super java/lang/Object 

.field public static i I


.method static <clinit>()V 
.limit stack 1
ldc 2
putstatic HelloWorld/i I
return 
.end method




.method public static main([Ljava/lang/String;)V
.limit stack 4
.limit locals 2
ldc 2
istore_2
ldc 1

ifeq ELSE_LABEL_0
ldc 3
iload 2
goto END_IF_ELSE_0
ELSE_LABEL_0:
END_IF_ELSE_0:
ldc 3
ldc 4
irem
iload 2
ldc 0
iload_2
invokestatic HelloWorld/a(IZ)I
return
.end method


.method public static a(ZI)I
.limit stack 3
.limit locals 3
ldc 2
istore_2
ldc 1

ifeq ELSE_LABEL_1
goto END_IF_ELSE_1
ELSE_LABEL_1:
END_IF_ELSE_1:
iload_2
ireturn
ldc 0
ireturn
.end method

