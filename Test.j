.class public Test
.super java/lang/Object 

.field public static i I

.method static <clinit>()V 
.limit stack 1
ldc 2
putstatic Test/i I
return 
.end method



.method public static add(II)I
.limit stack 3
.limit locals 2
iload_0
iload_1
iadd
getstatic Test/i I
iadd
ireturn
.end method

.method public static main([Ljava/lang/String;)V
.limit stack 6
.limit locals 2
ldc 1
ldc 2
getstatic Test/i I
imul
iadd
istore_2
getstatic Test/i I
iload_2
invokestatic Test/add(II)I
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc "Hello World"
invokevirtual java/io/PrintStream/println(Ljava/lang/String;)V
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc 1
invokevirtual java/io/PrintStream/println(I)V
return 
.end method

