import java.util.regex.Pattern;

public class EmailValidator implements UserValidator {

    private static final Pattern EMAIL_PATTERN =
        Pattern.compile(
            "^[_A-Za-z0-9-\\+]+(\\.[_A-Za-z0-9-]+)*"
            + "@[A-Za-z0-9-]+(\\.[A-Za-z0-9]+)*"
            + "(\\.[A-Za-z]{2,})$"
        );

    @Override
    public void validate(User user) {

        String email = user.getEmail();

        if (email == null || email.trim().isEmpty()) {
            throw new IllegalArgumentException(
                "Email should not be empty"
            );
        }

        if (!EMAIL_PATTERN.matcher(email).matches()) {
            throw new IllegalArgumentException(
                "Email is wrong format"
            );
        }
    }
}