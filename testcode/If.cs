namespace a {
    class b{
        static void Main(string[] args){
            if (a) {
                a();
            } else{
                b();
            }
        }
    }
}

// C# parser empty else without {}