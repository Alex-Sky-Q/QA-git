import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import java.time.Duration;
import static org.junit.jupiter.api.Assertions.*;

@DisplayName("Test BankAccount methods")
public class BankAccountTest {
    @Test
    @DisplayName("Deposit to a non-empty account")
    public void testDeposit() {
        double balance = 500;
        double depAmount = 500;
        BankAccount bankAccount = new BankAccount(balance, 0);
        bankAccount.deposit(depAmount);
        assertEquals(depAmount + balance, bankAccount.getBalance());
    }

    @Test
    @DisplayName("Withdraw to a min balance")
    public void testWithdrawToMin() {
        double balance = 500;
        double withdrawAmount = 500;
        BankAccount bankAccount = new BankAccount(balance, 0);
        bankAccount.withdraw(withdrawAmount);
        assertEquals(balance - withdrawAmount, bankAccount.getBalance());
    }

    @Test
    @DisplayName("Cannot withdraw more than minBalance")
    public void testWithdrawMoreThanMin() {
        double balance = 500;
        double withdrawAmount = 600;
        BankAccount bankAccount = new BankAccount(balance, 0);
        assertThrows(RuntimeException.class, () -> bankAccount.withdraw(withdrawAmount));
    }

    @Test
    @DisplayName("Account is active")
    public void testIsActive() {
        BankAccount bankAccount = new BankAccount(500, 0);
        assertTrue(bankAccount.isActive());
    }

    @Test
    @DisplayName("Holder name is not null")
    public void testHolderName() {
        BankAccount bankAccount = new BankAccount(500, 0);
        bankAccount.setHolderName("Travis");
        assertNotNull(bankAccount.getHolderName());
    }

    @Test
    @DisplayName("Withdraw is fast enough")
    public void testTimeOut() {
        BankAccount bankAccount = new BankAccount(500, 0);
        assertTimeout(Duration.ofMillis(10), () -> bankAccount.withdraw(500));
    }
}
