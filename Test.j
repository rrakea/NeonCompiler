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
.limit stack 4
.limit locals 2
ldc 1
ldc 2
getstatic Test/i i
imul
iadd
istore_1
getstatic Test/i i

iload_2

invokestatic Test/add(II)I
return 
.end method

