using System;
namespace  FibonacciRecursive
{
    class Program{
        static void Main (string []args){
            Console.WriteLine(Fib(5));
            Console.WriteLine(Fib(10));
        }

        static int Fib(int n){
            if (n > 0){
                if (n <= 2){
                    return 1;
                }
                return Fib (n - 1) + Fib(n - 2);
            }
            return 1;
        }
    }
}