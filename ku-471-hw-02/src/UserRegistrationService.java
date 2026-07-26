import java.util.List;

public class UserRegistrationService {

    private final List<UserValidator> validators;

    public UserRegistrationService(
        List<UserValidator> validators) {

        this.validators = validators;
    }

    public boolean register(User user) {

        for (UserValidator validator : validators) {
            validator.validate(user);
        }

        return true;
    }
}