.class public HelloWorld
.super java/lang/Object 

.field public static i I
.field public static x Z


.method static <clinit>()V 
.limit stack 3
ldc 2
putstatic HelloWorld/i I
ldc 1
ldc 0
iand
putstatic HelloWorld/x Z
return 
.end method




.method public static a(ZI)I
.limit stack 3
.limit locals 3
ldc 2
istore_2
ldc 1

ifeq ELSE_LABEL_0
goto END_IF_ELSE_0
ELSE_LABEL_0:
END_IF_ELSE_0:
iload_2
ireturn
ldc 0
ireturn
.end method


.method public static main([Ljava/lang/String;)V
.limit stack 5
.limit locals 2
ldc 2
istore_1
ldc 1

ifeq ELSE_LABEL_1
ldc 3
istore_1
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc "Inside Loop "
invokevirtual java/io/PrintStream/println(Ljava/lang/String;)V
getstatic java/lang/System/out Ljava/io/PrintStream;
iload_1
invokevirtual java/io/PrintStream/println(I)V
goto END_IF_ELSE_1
ELSE_LABEL_1:
END_IF_ELSE_1:
ldc 4
ldc 2
irem
istore_1
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc "HIde"
invokevirtual java/io/PrintStream/println(Ljava/lang/String;)V
getstatic java/lang/System/out Ljava/io/PrintStream;
iload_1
invokevirtual java/io/PrintStream/println(I)V
getstatic java/lang/System/out Ljava/io/PrintStream;
ldc "Consts: "
invokevirtual java/io/PrintStream/println(Ljava/lang/String;)V
getstatic java/lang/System/out Ljava/io/PrintStream;
getstatic HelloWorld/x Z
invokevirtual java/io/PrintStream/println(Z)V
getstatic java/lang/System/out Ljava/io/PrintStream;
getstatic HelloWorld/i I
invokevirtual java/io/PrintStream/println(I)V
return
.end method

