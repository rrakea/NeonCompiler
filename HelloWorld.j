.class public HelloWorld
.super java/lang/Object 

.field public static x Z
.field public static i I


.method static <clinit>()V 
.limit stack 3
ldc 1
ldc 0
iand
putstatic HelloWorld/x Z
ldc 2
putstatic HelloWorld/i I
return 
.end method




.method public static a(I)D
.limit stack 2
.limit locals 1
ldc2_w 0.0
dreturn
ldc 0.0
dreturn
.end method


.method public static main([Ljava/lang/String;)V
.limit stack 3
.limit locals 1
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc2_w 20
invokestatic HelloWorld/a(I)D
invokevirtual java/io/PrintStream/println(D)V
return
.end method

