import java.util.Random;

public class Die {
    private static final Random RANDOM = new Random();

    private int faceValue;

    public Die() {
        roll();
    }

    public int getFaceValue() {
        return faceValue;
    }

    private void setFaceValue(int faceValue) {
        this.faceValue = faceValue;
    }

    public void roll() {
        setFaceValue(RANDOM.nextInt(6) + 1);
    }
}
