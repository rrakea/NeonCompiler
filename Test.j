.class public Test
.super java/lang/Object 

.field public static i I

.method static <clinit>()V 
.limit locals 0
.limit stack 1
ldc 2
putstatic Test.j/i I
return 
.end method

.method public static add(II) I
.limit locals 0
.limit stack 2
ldc 0
istore 0
iload 3
iload 4
dadd
iload 2
dadd
return
.end method

.method public static main([Ljava/lang/String;) V
.limit locals 0
.limit stack 4
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

