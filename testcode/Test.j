.class public testcode/Test.cs
.super java/lang/Object 

.field public static i int

.method static <clinit>()V 
.limit locals 0
.limit stack 1
ldc 2
putstatic testcode/Test.j/i I
return 
.end method

.method public static main([Ljava/lang/String;) V
.limit stack 4
.limit locals 0
ldc 1
ldc 2
getstatic testcode/Test.cs/i I
imul
iadd
istore 0
iload 1
getstatic testcode/Test.cs/i I

invokestatic testcode/Test.cs/add()
.end method

