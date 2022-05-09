import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class BankAccountOrderedTest {
    static BankAccount bankAccount = new BankAccount(0, 0);

    @Test
    @Order(2)
    public void testWithdraw() {
        bankAccount.withdraw(500);
        assertEquals(0, bankAccount.getBalance());
    }

    @Test
    @Order(1)
    public void testDeposit() {
        bankAccount.deposit(500);
        assertEquals(500, bankAccount.getBalance());
    }
}
