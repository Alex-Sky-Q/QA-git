import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import java.time.Duration;
import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.*;

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
    @DisplayName("Balance can go below 0")
    public void testWithdrawBelowZero() {
        double balance = 500;
        double withdrawAmount = 600;
        BankAccount bankAccount = new BankAccount(balance, -1000);
        assertNotEquals(0, bankAccount.withdraw(withdrawAmount));
    }

    @Test
    @DisplayName("Cannot withdraw more than minBalance")
    public void testWithdrawMoreThanMin() {
        double balance = 500;
        double withdrawAmount = 600;
        BankAccount bankAccount = new BankAccount(balance, 0);
        assertThrows(RuntimeException.class, () -> bankAccount.withdraw(withdrawAmount),
                "Exception was not thrown");
    }

    @Test
    @DisplayName("Split account balance")
    public void testSplitAccount() {
        double balance = 10;
        int parts = 3;
        BankAccount bankAccount = new BankAccount(balance, 0);
        assertEquals(3.33, bankAccount.splitAccountBalance(parts), 0.01);
    }

    @Test
    @DisplayName("End-to-end process is successful")
    public void testEndToEnd() {
        double minBalance = 0;
        double depAmount = 500;
        double withdrawAmount = 300;
        BankAccount bankAccount = new BankAccount(0, minBalance);
        assumeTrue(bankAccount != null, "Account is null");
        assertAll(() -> bankAccount.setHolderName("Tracy"), ()-> bankAccount.deposit(depAmount),
                () -> bankAccount.withdraw(withdrawAmount));
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
        assumingThat(bankAccount.isActive(), () -> assertNotNull(bankAccount.getHolderName()));
    }

    @Test
    @DisplayName("Withdraw is fast enough")
    public void testTimeOut() {
        BankAccount bankAccount = new BankAccount(500, 0);
        assertTimeout(Duration.ofMillis(10), () -> bankAccount.withdraw(500));
    }
}
