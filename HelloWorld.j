.class public HelloWorld
.super java/lang/Object 

.field public static i I


.method static <clinit>()V 
.limit stack 1
ldc 2
putstatic HelloWorld/i I
return 
.end method




.method public static a(IZ)I
.limit stack 3
.limit locals 3
ldc 2
istore_2
ldc 1

ifeq ELSE_LABEL_0
goto END_IF_ELSE_0
ELSE_LABEL_0:
END_IF_ELSE_0:
iload_2
ireturn
ldc 0
ireturn
.end method


.method public static main([Ljava/lang/String;)V
.limit stack 4
.limit locals 2
ldc 2
istore_1
ldc 1

ifeq ELSE_LABEL_1
ldc 3
istore_1
goto END_IF_ELSE_1
ELSE_LABEL_1:
END_IF_ELSE_1:
ldc 3
ldc 4
irem
istore_1
ldc 0
iload_1
invokestatic HelloWorld/a(ZI)I
return
.end method

