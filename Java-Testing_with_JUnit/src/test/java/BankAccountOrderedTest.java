import org.junit.jupiter.api.*;
import org.junit.jupiter.api.condition.*;

import static org.junit.jupiter.api.Assertions.*;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class BankAccountOrderedTest {
    static BankAccount bankAccount;
    @BeforeAll
    public void testPrep() {
        bankAccount = new BankAccount(0, 0);
    }

    @BeforeEach
    public void testSetup() {
        System.out.println("Test setup");
    }

    @AfterEach
    public void testEnd() {
        System.out.println("Test teardown");
    }

    @DisabledOnJre({JRE.JAVA_18})
    @Test
    @Order(2)
    public void testWithdraw() {
        bankAccount.withdraw(500);
        assertEquals(0, bankAccount.getBalance());
    }

    @EnabledOnOs({OS.WINDOWS})
    @Test
    @Order(1)
    public void testDeposit() {
        bankAccount.deposit(500);
        assertEquals(500, bankAccount.getBalance());
    }
}
