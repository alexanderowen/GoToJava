
class loop{
    public static void main(String[] args) {
        int x = 0;
        for (;x < 10;) {
            x = x + 1;
        };
        int i = 0;
        for (; i < 10; i += 1) {
            System.out.println(i);
        }
        for (;;) {
            System.out.println("Infinite loop");
        }
    }
}
