using System;
namespace test {
    class test{
    static int i = 2;
        static void Main (string[] args) {
            int val = 1 + 2 * i;
            add(i, val);
        }

        static int add (int a, int b) {
            int c = 0;
            return a + b + c;
        }
    }
}

