.class public test
.super java/lang/Object 

.field public static i I


.method static <clinit>()V 
.limit stack 1
ldc 2
putstatic test/i I
return 
.end method




.method public static add(II)I
.limit stack 3
.limit locals 2
iload_0
iload_1
iadd
getstatic test/i I
iadd
ireturn
ldc 0
ireturn
.end method

