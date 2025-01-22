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
.limit locals 3
.limit stack 2
ldc 0
istore 1
iload 0
iload 1
dadd
iload 3
dadd
return
.end method

