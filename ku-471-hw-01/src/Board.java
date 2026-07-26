public class Board {
    private static final int SQUARE_COUNT = 40;

    private final Square[] squares;

    public Board() {
        squares = new Square[SQUARE_COUNT];

        for (int i = 0; i < squares.length; i++) {
            squares[i] = new Square("Square " + i);
        }
    }

    public Square getStartSquare() {
        return squares[0];
    }

    public Square getSquare(Square oldLocation, int faceValueTotal) {
        int oldIndex = findSquareIndex(oldLocation);
        int newIndex = (oldIndex + faceValueTotal) % squares.length;
        return squares[newIndex];
    }

    private int findSquareIndex(Square location) {
        for (int i = 0; i < squares.length; i++) {
            if (squares[i] == location) {
                return i;
            }
        }

        return 0;
    }
}
