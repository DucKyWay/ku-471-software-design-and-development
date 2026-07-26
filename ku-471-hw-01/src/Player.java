public class Player {
    private final String name;
    private final Piece piece;
    private final Board board;
    private final Die[] dice;

    public Player(String name, Piece piece, Board board, Die[] dice) {
        this.name = name;
        this.piece = piece;
        this.board = board;
        this.dice = dice;
    }

    public String getName() {
        return name;
    }

    public Piece getPiece() {
        return piece;
    }

    public void takeTurn() {
        int fv = 0;

        for (Die die : dice) {
            die.roll();
            fv += die.getFaceValue();
        }

        Square oldLocation = piece.getLocation();
        Square newLocation = board.getSquare(oldLocation, fv);
        piece.setLocation(newLocation);

        System.out.println(name + " roll " + fv + " move "
                + oldLocation.getName() + " to " + newLocation.getName());
    }
}
