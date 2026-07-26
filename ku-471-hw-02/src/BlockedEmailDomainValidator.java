import java.util.Set;

public class BlockedEmailDomainValidator
        implements UserValidator {

    private final Set<String> blockedDomains;

    public BlockedEmailDomainValidator( Set<String> blockedDomains) {

        this.blockedDomains = blockedDomains;
    }

    @Override
    public void validate(User user) {

        String email = user.getEmail();

        String domain = email.substring(
            email.indexOf("@") + 1
        );

        if (blockedDomains.contains(domain)) {
            throw new IllegalArgumentException(
                "Domain Email is not allow"
            );
        }
    }
}