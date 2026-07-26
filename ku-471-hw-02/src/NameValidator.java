import java.util.regex.Pattern;

public class NameValidator implements UserValidator {

    private static final Pattern NAME_PATTERN =
            Pattern.compile("[a-zA-Z]+");

    @Override
    public void validate(User user) {

        String name = user.getName();

        if (name == null || name.trim().isEmpty()) {
            throw new IllegalArgumentException(
                "Name should not empty"
            );
        }

        if (!NAME_PATTERN.matcher(name).matches()) {
            throw new IllegalArgumentException(
                "Name is wrong format"
            );
        }
    }
}