import org.junit.jupiter.api.*;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvFileSource;
import org.junit.jupiter.params.provider.ValueSource;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.*;

//@Execution(ExecutionMode.CONCURRENT)
@DisplayName("Test BankAccount methods")
@ExtendWith(BankAccountParamResolver.class)
public class BankAccountTest {
    @ParameterizedTest
    @CsvFileSource(resources = "name-amount.csv")
    @DisplayName("Deposit and set name")
    public void testDepositSetName(String name, Double depAmount, BankAccount bankAccount) {
        bankAccount.setHolderName(name);
        bankAccount.deposit(depAmount);
        assertEquals(name, bankAccount.getHolderName());
        assertEquals(depAmount, bankAccount.getBalance());
    }

    @ParameterizedTest
    @ValueSource(doubles = {100, 300.5, 0, -100})
    @DisplayName("Deposit to a non-empty account")
    public void testDeposit(double depAmount, BankAccount bankAccount) {
        bankAccount.deposit(depAmount);
        assertEquals(depAmount, bankAccount.getBalance());
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
    @Timeout(value = 1000, unit = TimeUnit.MILLISECONDS)
    @DisplayName("Account is active")
    public void testIsActive(BankAccount bankAccount) {
        assertTrue(bankAccount.isActive());
    }

    @Test
    @Disabled("Just for testing purpose")
    @DisplayName("Holder name is not null")
    public void testHolderName(BankAccount bankAccount) {
        bankAccount.setHolderName("Travis");
        assumingThat(bankAccount.isActive(), () -> assertNotNull(bankAccount.getHolderName()));
    }

    @Nested
    class PerformanceTest {
        @RepeatedTest(3)
        @DisplayName("Deposit is fast enough")
        public void testTimeOut(BankAccount bankAccount, RepetitionInfo repetitionInfo) {
            assertTimeout(Duration.ofMillis(10), () -> bankAccount.deposit(500));
            System.out.println(repetitionInfo.getCurrentRepetition() + " - " + bankAccount.getBalance());
        }
    }
}
