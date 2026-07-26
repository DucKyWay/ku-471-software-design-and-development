public class AgeValidator implements UserValidator {

    private final int minimumAge;

    public AgeValidator(int minimumAge) {
        this.minimumAge = minimumAge;
    }

    @Override
    public void validate(User user) {

        if (user.getAge() < minimumAge) {
            throw new IllegalArgumentException(
                "Age should be more than "
                    + minimumAge
                    + " years"
            );
        }
    }
}