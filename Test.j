.class public Test
.super java/lang/Object 

.field public static i int

.method static <clinit>()V 
.limit locals 0
.limit stack 1
ldc 2
putstatic Test.j/i I
return 
.end method

.method public static add(II) I
.limit stack 0
.limit locals 0
.end method

.method public static main([Ljava/lang/String;) V
.limit stack 4
.limit locals 0
ldc 1
ldc 2
getstatic Test/i I
imul
iadd
istore 0
getstatic Test/i I

iload 1
invokestatic Test/add(II)I
return 
.end method

