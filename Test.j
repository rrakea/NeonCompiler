.class public Test
.super java/lang/Object 

.field public static i I

.method static <clinit>()V 
.limit locals 0
.limit stack 1
ldc 2
putstatic Test/i I
return 
.end method

.method public static add(intint) I
.limit locals 3
.limit stack 2
ldc 0
istore 0
iload 0
iload 1
iadd
iload 2
iadd
ireturn
.end method

.method public static main(string[]) V
.limit locals 2
.limit stack 4
ldc 1
ldc 2
getstatic Test/i i
imul
iadd
istore 0
getstatic Test/i i

iload 1

invokestatic Test/add(II)I
return 
.end method

