public class MGame {
    private static final int DIE_COUNT = 2;

    private final int roundCount;
    private final Board board;
    private final Die[] dice;
    private final Player[] players;

    public MGame(int playerCount, int roundCount) {
        this.roundCount = roundCount;
        board = new Board();
        dice = createDice();
        players = createPlayers(playerCount);
    }

    public void playGame() {
        for (int round = 0; round < roundCount; round++) {
            System.out.println("Round: " + (round + 1));
            playRound();
        }
    }

    private void playRound() {
        for (Player player : players) {
            player.takeTurn();
        }
    }

    private Die[] createDice() {
        Die[] dice = new Die[DIE_COUNT];

        for (int i = 0; i < dice.length; i++) {
            dice[i] = new Die();
        }

        return dice;
    }

    private Player[] createPlayers(int playerCount) {
        Player[] players = new Player[playerCount];

        for (int i = 0; i < players.length; i++) {
            Piece piece = new Piece(board.getStartSquare());
            players[i] = new Player("Player " + (i + 1), piece, board, dice);
        }

        return players;
    }
}
