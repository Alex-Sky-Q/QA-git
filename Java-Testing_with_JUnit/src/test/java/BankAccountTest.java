import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
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
    public void testWithdraw() {
        double balance = 500;
        double withdrawAmount = 500;
        BankAccount bankAccount = new BankAccount(balance, 0);
        bankAccount.withdraw(withdrawAmount);
        assertEquals(balance - withdrawAmount, bankAccount.getBalance());
    }

    @Test
    @DisplayName("Account is active")
    public void testIsActive() {
        BankAccount bankAccount = new BankAccount(500, 0);
        assertTrue(bankAccount.isActive());
    }
}
