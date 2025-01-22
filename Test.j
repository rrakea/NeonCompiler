.class public Test
.super java/lang/Object 

.field public static i I

.method static <clinit>()V 
.limit stack 1
ldc 2
putstatic Test/i I
return 
.end method

.method public static main([Ljava/lang/String;)V
.limit stack 6
.limit locals 2
ldc 1
ldc 2
getstatic Test/i I
imul
iadd
istore_1
getstatic Test/i I

iload_1

invokestatic Test/add(II)I
return 
.end method

.method public static add(II)I
.limit stack 4
.limit locals 3
ldc 0
istore_2
iload_0
iload_1
iadd
iload_2
iadd
ireturn
.end method

