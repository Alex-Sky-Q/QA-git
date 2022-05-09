import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.EnumSource;
import java.time.DayOfWeek;
import static org.junit.jupiter.api.Assertions.*;

public class BankOfficeTest {
    @ParameterizedTest
    @EnumSource(value = DayOfWeek.class, names = {"SATURDAY", "SUNDAY"})
    @DisplayName("Day off is Sunday, not Saturday")
    public void testDaysOff(DayOfWeek day) {
        BankOffice bankOffice = new BankOffice();
        assertEquals(day.toString(), bankOffice.daysOff());
    }
}
