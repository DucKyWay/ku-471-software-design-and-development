import java.util.List;
import java.util.Set;

void main() {

    UserRegistrationService registrationService =
    new UserRegistrationService(
        List.of(
            new NameValidator(),
            new EmailValidator(),
            new BlockedEmailDomainValidator(
                Set.of(
                    "dom1.cc",
                    "dom2.cc",
                    "dom3.cc"
                )
            ),
            new AgeValidator(20)
        )
    );

    User user = new User(
        "test",
        "test@test.test",
        25
    );

    boolean result = registrationService.register(user);
    System.out.println(result);
}